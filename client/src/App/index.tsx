import React from 'react';
import './index.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import ChatRoom from './ChatRoom';

const App: React.FC = () => (
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<ChatRoom />} />
      {/* <Route path="/active" element={<TodoMVC />} />
      <Route path="/completed" element={<TodoMVC />} />
      <Route path="*" element={<NotFound />} /> */}
    </Routes>
</BrowserRouter>
)

export default App;
