import React, { useState } from 'react';
import axios from 'axios';
import './SearchPage.css';

const SearchPage = () => {
    const [query, setQuery] = useState('');
    const [suggestions, setSuggestions] = useState([]);
    const [searchResults, setSearchResults] = useState([]);

    const handleInputChange = (e) => {
        const value = e.target.value;
        setQuery(value);

        if (value.length > 2) {
            axios.get(`http://localhost:1710/search?q=${value}`)
                .then(response => {
                    setSuggestions(response.data);
                })
                .catch(error => {
                    console.error('There was an error fetching suggestions!', error);
                });
        } else {
            setSuggestions([]);
        }
    };

    const handleSearch = (e) => {
        e.preventDefault();
        axios.get(`http://localhost:1710/search?q=${query}`)
            .then(response => {
                setSearchResults(response.data);
            })
            .catch(error => {
                console.error('There was an error fetching search results!', error);
            });
    };

    const handleSuggestionClick = (suggestion) => {
        setQuery(suggestion.name);
        setSuggestions([]);
        setSearchResults([suggestion]);
    };

    return (
        <div className="search-page">
            <h1>Search Parts</h1>
            <form onSubmit={handleSearch}>
                <input
                    type="text"
                    value={query}
                    onChange={handleInputChange}
                    placeholder="Search for parts..."
                />
                <button type="submit">Search</button>
            </form>
            {suggestions.length > 0 && (
                <ul className="suggestions-list">
                    {suggestions.map((suggestion) => (
                        <li key={suggestion.id} onClick={() => handleSuggestionClick(suggestion)}>
                            {suggestion.name}
                        </li>
                    ))}
                </ul>
            )}
            <ul className="search-results">
                {searchResults.map((result) => (
                    <li key={result.id}>
                        <div className="part-details">
                            <strong>{result.name}</strong>
                            <div><strong>Price:</strong> ${result.price}</div>
                            <div><strong>SKU:</strong> {result.sku}</div>
                            <div><strong>Description:</strong> {result.description}</div>
                            <div><strong>Location:</strong> {result.location}</div>
                            <div><strong>Attributes:</strong> {result.attributes ? Object.entries(result.attributes).map(([key, value]) => `${key}: ${value}`).join(', ') : 'None'}</div>
                            <div><strong>Fitment Data:</strong> {result.fitmentData ? result.fitmentData.join(', ') : 'None'}</div>
                            <div><strong>Shipment:</strong> 
                                Weight: {result.shipment ? result.shipment.weight : 'N/A'}, 
                                Size: {result.shipment ? result.shipment.size : 'N/A'}, 
                                Hazardous: {result.shipment ? (result.shipment.hazardous ? 'Yes' : 'No') : 'N/A'}, 
                                Fragile: {result.shipment ? (result.shipment.fragile ? 'Yes' : 'No') : 'N/A'}
                            </div>
                            <div><strong>Metadata:</strong> {result.metadata ? Object.entries(result.metadata).map(([key, value]) => `${key}: ${value}`).join(', ') : 'None'}</div>
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default SearchPage;
