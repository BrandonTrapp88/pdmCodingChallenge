import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import './PartsList.css';

const PartsList = () => {
  const [parts, setParts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [expandedPart, setExpandedPart] = useState(null);
  const [version, setVersion] = useState(null);
  const [versionData, setVersionData] = useState(null);
  const [versionDetails, setVersionDetails] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    axios.get('http://localhost:1710/parts')
      .then(response => {
        setParts(response.data || []);
        setLoading(false);
      })
      .catch(error => {
        setError(error);
        setLoading(false);
      });
  }, []);

  const toggleExpand = (id) => {
    setExpandedPart(expandedPart === id ? null : id);
  };

  const deletePart = (id) => {
    axios.delete(`http://localhost:1710/parts/${id}`)
      .then(response => {
        setParts(parts.filter(part => part.id !== id));
      })
      .catch(error => {
        console.error('There was an error deleting the part!', error);
      });
  };

  const getPartVersions = (id) => {
    axios.get(`http://localhost:1710/parts/${id}/versions`)
      .then(response => {
        setVersionDetails(response.data);
      })
      .catch(error => {
        console.error('There was an error fetching the part versions!', error);
      });
  };

  const showVersion = (id, version) => {
    axios.get(`http://localhost:1710/parts/${id}/version/${version}`)
      .then(response => {
        const part = response.data;
        setVersionData(part);
      })
      .catch(error => {
        console.error('There was an error fetching the part version!', error);
      });
  };

  const editPart = (part) => {
    navigate(`/add?id=${part.id}`);
  };

  const getVersion = (id, version) => {
    if (version) {
      showVersion(id, version);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error fetching parts: {error.message}</div>;
  }

  return (
    <div className="parts-list-container">
      <h1>Vehicle Parts Inventory</h1>
      <ul className="parts-list">
        {parts.length === 0 ? (
          <li>No parts available</li>
        ) : (
          parts.map(part => (
            <li key={part.id} className="part-item">
              <div onClick={() => toggleExpand(part.id)} className="part-name">
                <strong>{part.name}</strong>
              </div>
              {expandedPart === part.id && (
                <div className="part-details">
                  <div><strong>Price:</strong> ${part.price}</div>
                  <div><strong>Images:</strong> {part.images ? part.images.join(', ') : 'None'}</div>
                  <div><strong>SKU:</strong> {part.sku}</div>
                  <div><strong>Description:</strong> {part.description}</div>
                  <div><strong>Attributes:</strong> {part.attributes ? Object.entries(part.attributes).map(([key, value]) => `${key}: ${value}`).join(', ') : 'None'}</div>
                  <div><strong>Fitment Data:</strong> {part.fitmentData ? part.fitmentData.join(', ') : 'None'}</div>
                  <div><strong>Location:</strong> {part.location}</div>
                  <div><strong>Shipment:</strong> 
                    Weight: {part.shipment ? part.shipment.weight : 'N/A'}, 
                    Size: {part.shipment ? part.shipment.size : 'N/A'}, 
                    Hazardous: {part.shipment ? (part.shipment.hazardous ? 'Yes' : 'No') : 'N/A'}, 
                    Fragile: {part.shipment ? (part.shipment.fragile ? 'Yes' : 'No') : 'N/A'}
                  </div>
                  <div><strong>Metadata:</strong> {versionData ? (versionData.metadata ? Object.entries(versionData.metadata).map(([key, value]) => `${key}: ${value}`).join(', ') : 'None') : (part.metadata ? Object.entries(part.metadata).map(([key, value]) => `${key}: ${value}`).join(', ') : 'None')}</div>
                  <button onClick={() => deletePart(part.id)} className="delete-button">Delete</button>
                  <button onClick={() => editPart(part)} className="edit-button">Edit</button>
                  <input type="number" value={version || ''} onChange={(e) => setVersion(e.target.value)} placeholder="Enter version" />
                  <button onClick={() => getVersion(part.id, version)} className="version-button">Get Version</button>
                </div>
              )}
            </li>
          ))
        )}
      </ul>
    </div>
  );
};

export default PartsList;
