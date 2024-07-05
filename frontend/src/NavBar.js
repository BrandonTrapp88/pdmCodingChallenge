import React from 'react';
import { Link } from 'react-router-dom';
import './NavBar.css';

const NavBar = () => {
  return (
    <nav>
      <ul>
        <li>
          <Link to="/add">Add Part</Link>
        </li>
        <li>
          <Link to="/">Parts List</Link>
        </li>
        <li>
          {/*<Link to="/search">Search Parts</Link> {/* Add search link */}
        </li>
      </ul>
    </nav>
  );
};

export default NavBar;
