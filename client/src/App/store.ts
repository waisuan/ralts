import { configureStore } from '@reduxjs/toolkit'
import chatMiddleware from './ChatRoom/chatMiddleware'
import chatReducer from './ChatRoom/chatSlice'

export const store = configureStore({
  reducer: {
    chat: chatReducer
  },
  middleware: (getDefaultMiddleware) => {
    return getDefaultMiddleware().concat([chatMiddleware])
  },
})

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch