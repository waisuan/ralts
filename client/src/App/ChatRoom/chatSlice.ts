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
  email: string,
  familyName: string,
  givenName: string,
  name: string,
  picture: string,
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
      state.isConnecting = true;
    }),
    setConnState: ((state, action) => {
      state.isConnecting = false;
      state.connected = action.payload;
    }),
    disconnect: ((state) => {
      state.userProfile = undefined;
      localStorage.removeItem('chat_sess_token');
    }),
    setUserProfile: ((state, action) => {
      const authToken = action.payload;
      state.userProfile = jwt_decode<UserProfile>(authToken);
      console.log(state.userProfile);
      localStorage.setItem('chat_sess_token', authToken);
    }),
    receiveMessage: ((state, action) => {
      console.log(action.payload);
      const payload = JSON.parse(action.payload);
      state.messages.push(payload);
    }),
    sendMessage: ((state, action) => {
      // state.messages.push(action.payload);
    }),
  },
})

export const chatSliceActions = chatSlice.actions

export const getMessages = (state: RootState) => state.chat.messages

export const getUserProfile = (state: RootState) => state.chat.userProfile

export default chatSlice.reducer