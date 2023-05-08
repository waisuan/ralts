import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'
import type { RootState } from '../store'
import jwt_decode from "jwt-decode"
import { DispatchProp } from 'react-redux'

interface Message {
  username: string,
  message: string,
  createdAt: string
}

interface UserProfile {
  username: string
}

interface ChatState {
  messages: Message[],
  connected: boolean,
  isConnecting: boolean,
  userProfile?: UserProfile,
  connCount: number
}

const initialState: ChatState = {
  messages: [],
  connected: false,
  isConnecting: false,
  connCount: 0
}

export const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    initConnection: ((state, action) => {
      state.messages = [];
      state.isConnecting = true;
    }),
    setConnState: ((state, action) => {
      state.isConnecting = false;
      state.connected = action.payload;
    }),
    disconnect: ((state) => {
      state.messages = [];
      state.userProfile = undefined;
      localStorage.removeItem('chat_sess_token');
    }),
    setUserProfile: ((state, action) => {
      state.userProfile = action.payload;
      localStorage.setItem('chat_sess_token', JSON.stringify(action.payload));
    }),
    receiveMessage: ((state, action) => {
      console.log(action.payload);
      const payload = JSON.parse(action.payload);
      state.messages.push(payload);
    }),
    sendMessage: ((state, action) => {
      //
    }),
    setConnCount: ((state, action) => {
      state.connCount = action.payload.count;
    }),
  },
})

export const chatSliceActions = chatSlice.actions

export const fetchConnCount = () => async (dispatch: any) => {
  await fetch('http://localhost:8001/conn_count')
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      dispatch(chatSliceActions.setConnCount(data))
    })
    .catch((err) => {
      console.log(err.message);
    });
  // dispatch(usersLoading());
  // const response = await usersAPI.fetchAll()
  // dispatch(usersReceived(response.data));
}

export const getMessages = (state: RootState) => state.chat.messages

export const getConnCount = (state: RootState) => state.chat.connCount

export const getUserProfile = (state: RootState) => state.chat.userProfile

export default chatSlice.reducer