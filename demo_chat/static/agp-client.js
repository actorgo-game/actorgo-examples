(function (root) {
    "use strict";

    var textEncoder = new TextEncoder();
    var textDecoder = new TextDecoder("utf-8");

    function concat(parts) {
        var size = parts.reduce(function (total, part) { return total + part.length; }, 0);
        var output = new Uint8Array(size);
        var offset = 0;
        parts.forEach(function (part) {
            output.set(part, offset);
            offset += part.length;
        });
        return output;
    }

    function bytes(value) {
        if (!value) return new Uint8Array(0);
        if (value instanceof Uint8Array) return value;
        if (value instanceof ArrayBuffer) return new Uint8Array(value);
        return new Uint8Array(value);
    }

    function varint(value) {
        var current = Number(value);
        var output = [];
        do {
            var next = current % 128;
            current = Math.floor(current / 128);
            if (current > 0) next += 128;
            output.push(next);
        } while (current > 0);
        return new Uint8Array(output);
    }

    function varintField(field, value) {
        return concat([varint(field * 8), varint(value)]);
    }

    function bytesField(field, value) {
        var data = bytes(value);
        return concat([varint(field * 8 + 2), varint(data.length), data]);
    }

    function Reader(value) {
        this.data = bytes(value);
        this.offset = 0;
    }

    Reader.prototype.done = function () { return this.offset >= this.data.length; };
    Reader.prototype.varint = function () {
        var result = 0;
        var multiplier = 1;
        while (this.offset < this.data.length && multiplier <= Math.pow(2, 63)) {
            var current = this.data[this.offset++];
            result += (current & 127) * multiplier;
            if ((current & 128) === 0) return result;
            multiplier *= 128;
        }
        throw new Error("无效的 Protobuf varint");
    };
    Reader.prototype.block = function () {
        var length = this.varint();
        var end = this.offset + length;
        if (end > this.data.length) throw new Error("Protobuf 字段越界");
        var result = this.data.slice(this.offset, end);
        this.offset = end;
        return result;
    };
    Reader.prototype.skip = function (wireType) {
        if (wireType === 0) this.varint();
        else if (wireType === 1) this.offset += 8;
        else if (wireType === 2) this.offset += this.varint();
        else if (wireType === 5) this.offset += 4;
        else throw new Error("不支持的 Protobuf wire type: " + wireType);
        if (this.offset > this.data.length) throw new Error("Protobuf 字段越界");
    };

    function requestPacket(requestID, methodID, timeout, body, codec) {
        var request = concat([
            varintField(1, requestID),
            varintField(2, methodID),
            varintField(3, timeout),
            bytesField(4, body)
        ]);
        return concat([bytesField(1, request), varintField(5, codec)]);
    }

    function notifyPacket(methodID, body, codec) {
        var notify = concat([varintField(1, methodID), bytesField(2, body)]);
        return concat([bytesField(3, notify), varintField(5, codec)]);
    }

    function decodeResponse(value) {
        var reader = new Reader(value);
        var result = {requestID: 0, code: 0, message: "", body: new Uint8Array(0)};
        while (!reader.done()) {
            var key = reader.varint();
            var field = Math.floor(key / 8);
            var wire = key & 7;
            if (field === 1 && wire === 0) result.requestID = reader.varint();
            else if (field === 2 && wire === 0) result.code = reader.varint();
            else if (field === 3 && wire === 2) result.message = textDecoder.decode(reader.block());
            else if (field === 4 && wire === 2) result.body = reader.block();
            else reader.skip(wire);
        }
        return result;
    }

    function decodeNotify(value) {
        var reader = new Reader(value);
        var result = {methodID: 0, body: new Uint8Array(0)};
        while (!reader.done()) {
            var key = reader.varint();
            var field = Math.floor(key / 8);
            var wire = key & 7;
            if (field === 1 && wire === 0) result.methodID = reader.varint();
            else if (field === 2 && wire === 2) result.body = reader.block();
            else reader.skip(wire);
        }
        return result;
    }

    function decodePacket(value) {
        var reader = new Reader(value);
        var packet = {codec: 0, response: null, notify: null};
        while (!reader.done()) {
            var key = reader.varint();
            var field = Math.floor(key / 8);
            var wire = key & 7;
            if (field === 2 && wire === 2) packet.response = decodeResponse(reader.block());
            else if (field === 3 && wire === 2) packet.notify = decodeNotify(reader.block());
            else if (field === 5 && wire === 0) packet.codec = reader.varint();
            else reader.skip(wire);
        }
        return packet;
    }

    function decodeHandshake(value) {
        var reader = new Reader(value);
        var result = {version: 0, heartbeat: 0, connectionID: ""};
        while (!reader.done()) {
            var key = reader.varint();
            var field = Math.floor(key / 8);
            var wire = key & 7;
            if (field === 1 && wire === 0) result.version = reader.varint();
            else if (field === 2 && wire === 0) result.heartbeat = reader.varint();
            else if (field === 4 && wire === 2) result.connectionID = textDecoder.decode(reader.block());
            else reader.skip(wire);
        }
        return result;
    }

    function AGPError(code, message) {
        this.name = "AGPError";
        this.code = code;
        this.message = message || ("AGP 请求失败，code=" + code);
    }
    AGPError.prototype = Object.create(Error.prototype);

    function AGPJSONClient() {
        this.socket = null;
        this.ready = false;
        this.sequence = 0;
        this.pending = new Map();
        this.heartbeatTimer = null;
        this.onNotify = null;
    }

    AGPJSONClient.prototype.connect = function (host, port) {
        var self = this;
        self.disconnect();
        var scheme = root.location.protocol === "https:" ? "wss" : "ws";
        var endpoint = scheme + "://" + host + ":" + port + "/";

        return new Promise(function (resolve, reject) {
            var settled = false;
            var socket = new WebSocket(endpoint, "agp.v1");
            self.socket = socket;
            socket.binaryType = "arraybuffer";
            socket.onopen = function () {
                self._request(1, bytesField(1, new Uint8Array([1])), 1, false, 10000)
                    .then(function (body) {
                        var handshake = decodeHandshake(body);
                        if (handshake.version !== 1) throw new Error("服务端未协商 AGP/1");
                        self.ready = true;
                        self.startHeartbeat(handshake.heartbeat);
                        settled = true;
                        resolve(handshake);
                    }).catch(function (error) {
                        settled = true;
                        self.disconnect();
                        reject(error);
                    });
            };
            socket.onmessage = function (event) {
                if (event.data instanceof Blob) {
                    event.data.arrayBuffer().then(function (buffer) { self.handle(buffer); });
                } else {
                    self.handle(event.data);
                }
            };
            socket.onerror = function () {
                if (!settled) {
                    settled = true;
                    reject(new Error("无法连接网关 " + endpoint));
                }
            };
            socket.onclose = function () {
                self.ready = false;
                self.stopHeartbeat();
                self.failAll(new Error("聊天室连接已关闭"));
                if (!settled) {
                    settled = true;
                    reject(new Error("聊天室连接在握手前关闭"));
                }
            };
        });
    };

    AGPJSONClient.prototype.request = function (methodID, value, timeout) {
        if (!this.ready) return Promise.reject(new Error("聊天室尚未连接"));
        return this._request(methodID, textEncoder.encode(JSON.stringify(value || {})), 2, true, timeout || 10000);
    };

    AGPJSONClient.prototype.notify = function (methodID, value) {
        if (!this.ready) throw new Error("聊天室尚未连接");
        this.socket.send(notifyPacket(methodID, textEncoder.encode(JSON.stringify(value || {})), 2));
    };

    AGPJSONClient.prototype._request = function (methodID, body, codec, decodeJSON, timeout) {
        var self = this;
        if (!self.socket || self.socket.readyState !== WebSocket.OPEN) {
            return Promise.reject(new Error("WebSocket 未连接"));
        }
        self.sequence = self.sequence % 0xffffffff + 1;
        var requestID = self.sequence;
        return new Promise(function (resolve, reject) {
            var timer = setTimeout(function () {
                self.pending.delete(requestID);
                reject(new Error("请求超时，methodID=" + methodID));
            }, timeout + 1000);
            self.pending.set(requestID, {resolve: resolve, reject: reject, timer: timer, json: decodeJSON});
            self.socket.send(requestPacket(requestID, methodID, timeout, body, codec));
        });
    };

    AGPJSONClient.prototype.handle = function (value) {
        var packet = decodePacket(value);
        if (packet.response) {
            var pending = this.pending.get(packet.response.requestID);
            if (!pending) return;
            clearTimeout(pending.timer);
            this.pending.delete(packet.response.requestID);
            if (packet.response.code !== 0) {
                pending.reject(new AGPError(packet.response.code, packet.response.message));
            } else if (pending.json && packet.response.body.length > 0) {
                pending.resolve(JSON.parse(textDecoder.decode(packet.response.body)));
            } else {
                pending.resolve(packet.response.body);
            }
            return;
        }
        if (packet.notify && typeof this.onNotify === "function") {
            var body = packet.codec === 2 && packet.notify.body.length > 0
                ? JSON.parse(textDecoder.decode(packet.notify.body))
                : packet.notify.body;
            this.onNotify(packet.notify.methodID, body);
        }
    };

    AGPJSONClient.prototype.startHeartbeat = function (interval) {
        var self = this;
        self.stopHeartbeat();
        if (!interval) return;
        self.heartbeatTimer = setInterval(function () {
            self._request(2, varintField(1, Date.now()), 1, false, 10000)
                .catch(function () { self.disconnect(); });
        }, interval);
    };

    AGPJSONClient.prototype.stopHeartbeat = function () {
        if (this.heartbeatTimer) clearInterval(this.heartbeatTimer);
        this.heartbeatTimer = null;
    };

    AGPJSONClient.prototype.failAll = function (error) {
        this.pending.forEach(function (pending) {
            clearTimeout(pending.timer);
            pending.reject(error);
        });
        this.pending.clear();
    };

    AGPJSONClient.prototype.disconnect = function () {
        this.ready = false;
        this.stopHeartbeat();
        this.failAll(new Error("聊天室客户端已断开"));
        if (this.socket) {
            this.socket.onclose = null;
            this.socket.close();
            this.socket = null;
        }
    };

    root.ActorGoAGPJSONClient = AGPJSONClient;
})(window);
