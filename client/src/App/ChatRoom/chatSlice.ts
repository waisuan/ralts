import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'
import type { RootState } from '../store'
import jwt_decode from "jwt-decode"

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
}

const initialState: ChatState = {
  messages: [],
  connected: false,
  isConnecting: false,
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
  },
})

export const chatSliceActions = chatSlice.actions

export const getMessages = (state: RootState) => state.chat.messages

export const getUserProfile = (state: RootState) => state.chat.userProfile

export default chatSlice.reducer