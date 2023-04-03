
import { Middleware } from 'redux'
import { chatSliceActions } from './chatSlice';
 
const WS_URL = 'ws://localhost:8001/ws';

const chatMiddleware: Middleware = store => {
  let socket: WebSocket;

  return next => action => {
    if (chatSliceActions.initConnection.match(action)) { 
      const authToken = action.payload
      socket = new WebSocket(WS_URL + "?authorization=" + authToken)

      socket.onopen = () => {
        store.dispatch(chatSliceActions.setConnState(true));
        store.dispatch(chatSliceActions.setUserProfile(authToken))
      };

      socket.onerror = () => {
        store.dispatch(chatSliceActions.setConnState(false));
      };

      // socket.onclose

      socket.onmessage = (e) => {
        store.dispatch(chatSliceActions.receiveMessage(e.data));
      };
    } else if(chatSliceActions.sendMessage.match(action)) {
      const payload = JSON.stringify({
        userId: Math.floor(Math.random() * 2) === 0 ? "EvanSia" : "Rando",
        message: action.payload
      });
      socket.send(payload);
    }
 
    next(action);
  }
}
 
export default chatMiddleware;