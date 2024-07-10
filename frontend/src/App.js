import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import AddPartForm from './AddPartForm';
import PartsList from './PartsList';
import SearchPage from './SearchPage';
import UpdatePartForm from './UpdatePartForm';
import NavBar from './NavBar'; // Ensure NavBar is imported

const App = () => {
  return (
    <Router>
      <div>
        <NavBar /> {}
        <Routes>
          <Route path="/" element={<PartsList />} />
          <Route path="/add" element={<AddPartForm />} />
          <Route path="/update" element={<UpdatePartForm />} />
          <Route path="/search" element={<SearchPage />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;
