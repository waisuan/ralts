import React from 'react';
import './index.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import ChatRoom from './ChatRoom';
import NewsFeed from './NewsFeed';

const App: React.FC = () => (
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<ChatRoom />} />
      <Route path="/feed" element={<NewsFeed />} />
      {/* <Route path="/completed" element={<TodoMVC />} />
      <Route path="*" element={<NotFound />} /> */}
    </Routes>
</BrowserRouter>
)

export default App;
