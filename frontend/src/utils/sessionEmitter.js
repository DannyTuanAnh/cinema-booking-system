// Event emitter cho session expired
class SessionEventEmitter {
  constructor() {
    this.listeners = [];
  }

  subscribe(callback) {
    this.listeners.push(callback);
    return () => {
      this.listeners = this.listeners.filter((cb) => cb !== callback);
    };
  }

  emit() {
    this.listeners.forEach((callback) => callback());
  }
}

export const sessionExpiredEmitter = new SessionEventEmitter();
