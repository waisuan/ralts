
import { Middleware } from 'redux'
import { chatSliceActions } from './chatSlice';
 
const WS_URL = 'ws://localhost:8001/ws';

const chatMiddleware: Middleware = store => {
  let socket: WebSocket;

  return next => action => {
    if (chatSliceActions.initConnection.match(action)) {
      socket = new WebSocket(WS_URL)

      socket.onopen = () => {
        store.dispatch(chatSliceActions.setConnState(true))
        store.dispatch(chatSliceActions.setUserProfile(action.payload))
      };

      socket.onerror = () => {
        store.dispatch(chatSliceActions.setConnState(false))
      };

      // socket.onclose

      socket.onmessage = (e) => {
        store.dispatch(chatSliceActions.receiveMessage(e.data))
      };
    } else if(chatSliceActions.sendMessage.match(action)) {
      const state = store.getState().chat

      socket.send(JSON.stringify({
        userId: state.userProfile!!.username, //Math.floor(Math.random() * 2) === 0 ? "EvanSia" : "Rando",
        message: action.payload
      }));
    } else if (chatSliceActions.disconnect.match(action)) {
      if (socket !== undefined) {
        socket.close()        
      }
    }
 
    next(action)
  }
}
 
export default chatMiddleware;