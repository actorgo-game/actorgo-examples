(function (root) {
    "use strict";

    var CODEC_PROTOBUF = 1;
    var METHOD_HANDSHAKE = 1;
    var METHOD_HEARTBEAT = 2;

    function concatBytes(parts) {
        var length = 0;
        var i;
        for (i = 0; i < parts.length; i++) {
            length += parts[i].length;
        }

        var output = new Uint8Array(length);
        var offset = 0;
        for (i = 0; i < parts.length; i++) {
            output.set(parts[i], offset);
            offset += parts[i].length;
        }
        return output;
    }

    function asBytes(value) {
        if (!value) {
            return new Uint8Array(0);
        }
        if (value instanceof Uint8Array) {
            return value;
        }
        if (value instanceof ArrayBuffer) {
            return new Uint8Array(value);
        }
        return new Uint8Array(value);
    }

    function encodeVarint(value) {
        var current = BigInt(value);
        if (current < 0) {
            current = BigInt.asUintN(64, current);
        }

        var bytes = [];
        do {
            var next = Number(current & 0x7fn);
            current >>= 7n;
            if (current !== 0n) {
                next |= 0x80;
            }
            bytes.push(next);
        } while (current !== 0n);
        return new Uint8Array(bytes);
    }

    function varintField(fieldNumber, value) {
        return concatBytes([
            encodeVarint((fieldNumber << 3) | 0),
            encodeVarint(value)
        ]);
    }

    function bytesField(fieldNumber, value) {
        var bytes = asBytes(value);
        return concatBytes([
            encodeVarint((fieldNumber << 3) | 2),
            encodeVarint(bytes.length),
            bytes
        ]);
    }

    function Reader(bytes) {
        this.bytes = asBytes(bytes);
        this.offset = 0;
    }

    Reader.prototype.done = function () {
        return this.offset >= this.bytes.length;
    };

    Reader.prototype.readVarint = function () {
        var result = 0n;
        var shift = 0n;
        while (this.offset < this.bytes.length && shift < 70n) {
            var value = this.bytes[this.offset++];
            result |= BigInt(value & 0x7f) << shift;
            if ((value & 0x80) === 0) {
                return result;
            }
            shift += 7n;
        }
        throw new Error("无效的 Protobuf varint");
    };

    Reader.prototype.readLengthDelimited = function () {
        var length = Number(this.readVarint());
        var end = this.offset + length;
        if (length < 0 || end > this.bytes.length) {
            throw new Error("无效的 Protobuf 字段长度");
        }
        var value = this.bytes.slice(this.offset, end);
        this.offset = end;
        return value;
    };

    Reader.prototype.skip = function (wireType) {
        if (wireType === 0) {
            this.readVarint();
            return;
        }
        if (wireType === 1) {
            this.offset += 8;
        } else if (wireType === 2) {
            this.offset += Number(this.readVarint());
        } else if (wireType === 5) {
            this.offset += 4;
        } else {
            throw new Error("不支持的 Protobuf wire type: " + wireType);
        }
        if (this.offset > this.bytes.length) {
            throw new Error("Protobuf 字段超出消息边界");
        }
    };

    function decodeString(bytes) {
        return new TextDecoder("utf-8").decode(bytes);
    }

    function encodeRequestPacket(requestId, methodId, timeoutMs, body) {
        var request = concatBytes([
            varintField(1, requestId),
            varintField(2, methodId),
            varintField(3, timeoutMs),
            bytesField(4, body)
        ]);
        return concatBytes([
            bytesField(1, request),
            varintField(5, CODEC_PROTOBUF)
        ]);
    }

    function encodeHandshakeRequest() {
        return bytesField(1, new Uint8Array([1]));
    }

    function encodeHeartbeatRequest() {
        return varintField(1, Date.now());
    }

    function decodeResponse(bytes) {
        var reader = new Reader(bytes);
        var response = {requestId: 0, code: 0, message: "", body: new Uint8Array(0)};
        while (!reader.done()) {
            var key = Number(reader.readVarint());
            var fieldNumber = key >>> 3;
            var wireType = key & 7;
            if (fieldNumber === 1 && wireType === 0) {
                response.requestId = Number(reader.readVarint());
            } else if (fieldNumber === 2 && wireType === 0) {
                response.code = Number(reader.readVarint());
            } else if (fieldNumber === 3 && wireType === 2) {
                response.message = decodeString(reader.readLengthDelimited());
            } else if (fieldNumber === 4 && wireType === 2) {
                response.body = reader.readLengthDelimited();
            } else {
                reader.skip(wireType);
            }
        }
        return response;
    }

    function decodeNotify(bytes) {
        var reader = new Reader(bytes);
        var notify = {methodId: 0, body: new Uint8Array(0)};
        while (!reader.done()) {
            var key = Number(reader.readVarint());
            var fieldNumber = key >>> 3;
            var wireType = key & 7;
            if (fieldNumber === 1 && wireType === 0) {
                notify.methodId = Number(reader.readVarint());
            } else if (fieldNumber === 2 && wireType === 2) {
                notify.body = reader.readLengthDelimited();
            } else {
                reader.skip(wireType);
            }
        }
        return notify;
    }

    function decodePacket(bytes) {
        var reader = new Reader(bytes);
        var packet = {codec: 0, response: null, notify: null};
        while (!reader.done()) {
            var key = Number(reader.readVarint());
            var fieldNumber = key >>> 3;
            var wireType = key & 7;
            if (fieldNumber === 2 && wireType === 2) {
                packet.response = decodeResponse(reader.readLengthDelimited());
            } else if (fieldNumber === 3 && wireType === 2) {
                packet.notify = decodeNotify(reader.readLengthDelimited());
            } else if (fieldNumber === 5 && wireType === 0) {
                packet.codec = Number(reader.readVarint());
            } else {
                reader.skip(wireType);
            }
        }
        return packet;
    }

    function decodeHandshakeResponse(bytes) {
        var reader = new Reader(bytes);
        var response = {
            protocolVersion: 0,
            heartbeatIntervalMs: 0,
            maxFrameSize: 0,
            connectionId: ""
        };
        while (!reader.done()) {
            var key = Number(reader.readVarint());
            var fieldNumber = key >>> 3;
            var wireType = key & 7;
            if (fieldNumber === 1 && wireType === 0) {
                response.protocolVersion = Number(reader.readVarint());
            } else if (fieldNumber === 2 && wireType === 0) {
                response.heartbeatIntervalMs = Number(reader.readVarint());
            } else if (fieldNumber === 3 && wireType === 0) {
                response.maxFrameSize = Number(reader.readVarint());
            } else if (fieldNumber === 4 && wireType === 2) {
                response.connectionId = decodeString(reader.readLengthDelimited());
            } else {
                reader.skip(wireType);
            }
        }
        return response;
    }

    function AGPError(code, message) {
        this.name = "AGPError";
        this.code = code;
        this.message = message || ("AGP 请求失败，code=" + code);
        if (Error.captureStackTrace) {
            Error.captureStackTrace(this, AGPError);
        }
    }
    AGPError.prototype = Object.create(Error.prototype);
    AGPError.prototype.constructor = AGPError;

    function AGPClient() {
        this.socket = null;
        this.ready = false;
        this.requestId = 0;
        this.pending = new Map();
        this.heartbeatTimer = null;
        this.connectionId = "";
        this.onNotify = null;
    }

    AGPClient.prototype.connect = function (host, port, path) {
        var self = this;
        self.disconnect();

        var scheme = root.location && root.location.protocol === "https:" ? "wss" : "ws";
        var endpoint = scheme + "://" + host + ":" + port + (path || "/");

        return new Promise(function (resolve, reject) {
            var settled = false;
            var socket;
            try {
                socket = new WebSocket(endpoint, "agp.v1");
            } catch (error) {
                reject(error);
                return;
            }

            self.socket = socket;
            socket.binaryType = "arraybuffer";

            socket.onopen = function () {
                self._request(METHOD_HANDSHAKE, encodeHandshakeRequest(), 10000, true)
                    .then(function (body) {
                        var handshake = decodeHandshakeResponse(body);
                        if (handshake.protocolVersion !== 1) {
                            throw new Error("服务端未协商 AGP/1");
                        }
                        self.ready = true;
                        self.connectionId = handshake.connectionId;
                        self._startHeartbeat(handshake.heartbeatIntervalMs);
                        settled = true;
                        resolve(handshake);
                    })
                    .catch(function (error) {
                        settled = true;
                        self.disconnect();
                        reject(error);
                    });
            };

            socket.onmessage = function (event) {
                if (event.data instanceof Blob) {
                    event.data.arrayBuffer().then(function (buffer) {
                        self._handleMessage(new Uint8Array(buffer));
                    }).catch(function (error) {
                        self._failAll(error);
                    });
                    return;
                }
                self._handleMessage(new Uint8Array(event.data));
            };

            socket.onerror = function () {
                if (!settled) {
                    settled = true;
                    reject(new Error("无法连接网关 " + endpoint));
                }
            };

            socket.onclose = function () {
                self.ready = false;
                self._stopHeartbeat();
                self._failAll(new Error("网关连接已关闭"));
                if (!settled) {
                    settled = true;
                    reject(new Error("网关连接在握手前关闭"));
                }
            };
        });
    };

    AGPClient.prototype.request = function (methodId, body, timeoutMs) {
        if (!this.ready) {
            return Promise.reject(new Error("AGP 网关尚未连接"));
        }
        return this._request(methodId, body, timeoutMs || 10000, false);
    };

    AGPClient.prototype._request = function (methodId, body, timeoutMs, allowBeforeReady) {
        var self = this;
        if (!self.socket || self.socket.readyState !== WebSocket.OPEN || (!allowBeforeReady && !self.ready)) {
            return Promise.reject(new Error("AGP WebSocket 未连接"));
        }

        self.requestId = (self.requestId % 0xffffffff) + 1;
        var requestId = self.requestId;
        var packet = encodeRequestPacket(requestId, methodId, timeoutMs, body);

        return new Promise(function (resolve, reject) {
            var timer = setTimeout(function () {
                self.pending.delete(requestId);
                reject(new Error("AGP 请求超时，methodId=" + methodId));
            }, timeoutMs + 1000);

            self.pending.set(requestId, {resolve: resolve, reject: reject, timer: timer});
            try {
                self.socket.send(packet);
            } catch (error) {
                clearTimeout(timer);
                self.pending.delete(requestId);
                reject(error);
            }
        });
    };

    AGPClient.prototype._handleMessage = function (bytes) {
        var packet;
        try {
            packet = decodePacket(bytes);
        } catch (error) {
            this._failAll(error);
            this.disconnect();
            return;
        }

        if (packet.response) {
            var pending = this.pending.get(packet.response.requestId);
            if (!pending) {
                return;
            }
            clearTimeout(pending.timer);
            this.pending.delete(packet.response.requestId);
            if (packet.response.code !== 0) {
                pending.reject(new AGPError(packet.response.code, packet.response.message));
            } else {
                pending.resolve(packet.response.body);
            }
            return;
        }

        if (packet.notify && typeof this.onNotify === "function") {
            this.onNotify(packet.notify);
        }
    };

    AGPClient.prototype._startHeartbeat = function (intervalMs) {
        var self = this;
        self._stopHeartbeat();
        if (!intervalMs) {
            return;
        }
        self.heartbeatTimer = setInterval(function () {
            if (!self.ready) {
                return;
            }
            self._request(METHOD_HEARTBEAT, encodeHeartbeatRequest(), 10000, false)
                .catch(function (error) {
                    console.error("AGP heartbeat failed", error);
                    self.disconnect();
                });
        }, intervalMs);
    };

    AGPClient.prototype._stopHeartbeat = function () {
        if (this.heartbeatTimer) {
            clearInterval(this.heartbeatTimer);
            this.heartbeatTimer = null;
        }
    };

    AGPClient.prototype._failAll = function (error) {
        this.pending.forEach(function (pending) {
            clearTimeout(pending.timer);
            pending.reject(error);
        });
        this.pending.clear();
    };

    AGPClient.prototype.disconnect = function () {
        this.ready = false;
        this._stopHeartbeat();
        this._failAll(new Error("AGP 客户端已断开"));
        if (this.socket) {
            this.socket.onclose = null;
            this.socket.close();
            this.socket = null;
        }
    };

    root.ActorGoAGPClient = AGPClient;
    root.ActorGoAGPError = AGPError;
})(typeof window !== "undefined" ? window : globalThis);
