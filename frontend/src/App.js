import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import AddPartForm from './AddPartForm';
import PartsList from './PartsList';
import NavBar from './NavBar';
import SearchPage from './SearchPage';

const App = () => {
  return (
    <Router>
      <NavBar />
      <Routes>
        <Route path="/add" element={<AddPartForm />} />
        <Route path="/" element={<PartsList />} />
        <Route path="/search" element={<SearchPage/> }/> {/* Add search route */}
      </Routes>
    </Router>
  );
};

export default App;
