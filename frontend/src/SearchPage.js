import React, { useState } from 'react';
import axios from 'axios';

const SearchPage = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [results, setResults] = useState([]);

  const handleSearch = () => {
    axios.get(`http://localhost:1710/parts/search?name=${searchTerm}`)
      .then(response => {
        setResults(response.data);
      })
      .catch(error => {
        console.error('There was an error searching for parts!', error);
      });
  };

  return (
    <div>
      <h1>Search Parts</h1>
      <div>
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Enter part name"
        />
        <button onClick={handleSearch}>Search</button>
      </div>
      <div>
        {results.length > 0 ? (
          <ul>
            {results.map(part => (
              <li key={part.id}>
                <strong>{part.name}</strong> - ${part.price}
              </li>
            ))}
          </ul>
        ) : (
          <p>No results found</p>
        )}
      </div>
    </div>
  );
};

export default SearchPage;
