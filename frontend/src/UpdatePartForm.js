import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useLocation, useNavigate } from 'react-router-dom';
import './UpdatePartForm.css';

const UpdatePartForm = () => {
  const [name, setName] = useState('');
  const [price, setPrice] = useState('');
  const [images, setImages] = useState('');
  const [sku, setSku] = useState('');
  const [description, setDescription] = useState('');
  const [attributes, setAttributes] = useState('');
  const [fitmentData, setFitmentData] = useState('');
  const [location, setLocation] = useState('');
  const [shipmentWeight, setShipmentWeight] = useState('');
  const [shipmentSize, setShipmentSize] = useState('');
  const [shipmentHazardous, setShipmentHazardous] = useState(false);
  const [shipmentFragile, setShipmentFragile] = useState(false);
  const [metadata, setMetadata] = useState('');
  const [partId, setPartId] = useState(null);

  const locationData = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const queryParams = new URLSearchParams(locationData.search);
    const id = queryParams.get('id');
    if (id) {
      setPartId(id);
      axios.get(`http://localhost:1710/parts/${id}`)
        .then(response => {
          const part = response.data;
          setName(part.name);
          setPrice(part.price);
          setImages(part.images.join(','));
          setSku(part.sku);
          setDescription(part.description);
          setAttributes(Object.entries(part.attributes).map(([key, value]) => `${key}:${value}`).join(','));
          setFitmentData(part.fitmentData.join(','));
          setLocation(part.location);
          setShipmentWeight(part.shipment.weight);
          setShipmentSize(part.shipment.size);
          setShipmentHazardous(part.shipment.hazardous);
          setShipmentFragile(part.shipment.fragile);
          setMetadata(Object.entries(part.metadata).map(([key, value]) => `${key}:${value}`).join(','));
        })
        .catch(error => {
          console.error('There was an error fetching the part!', error);
        });
    }
  }, [locationData]);

  const handleSubmit = (event) => {
    event.preventDefault();

    const part = {
      name,
      price: parseFloat(price),
      images: images ? images.split(',') : [],
      sku,
      description,
      attributes: attributes ? attributes.split(',').reduce((acc, attr) => {
        const [key, value] = attr.split(':');
        acc[key.trim()] = value.trim();
        return acc;
      }, {}) : {},
      fitmentData: fitmentData ? fitmentData.split(',') : [],
      location,
      shipment: {
        weight: shipmentWeight ? parseFloat(shipmentWeight) : 0,
        size: shipmentSize,
        hazardous: shipmentHazardous,
        fragile: shipmentFragile
      },
      metadata: metadata ? metadata.split(',').reduce((acc, meta) => {
        const [key, value] = meta.split(':');
        acc[key.trim()] = value.trim();
        return acc;
      }, {}) : {}
    };

    console.log('Updating part:', part);

    axios.patch(`http://localhost:1710/parts/${partId}`, part)
      .then(response => {
        console.log('Part updated successfully:', response.data);
        alert('Part updated successfully!');
        navigate('/');
      })
      .catch(error => {
        console.error('There was an error updating the part!', error);
      });
  };

  return (
    <form onSubmit={handleSubmit} className="update-part-form">
      <div className="form-group">
        <label>Part Name:</label>
        <input type="text" value={name} onChange={(e) => setName(e.target.value)} required />
      </div>
      <div className="form-group">
        <label>Price:</label>
        <input type="number" value={price} onChange={(e) => setPrice(e.target.value)} required />
      </div>
      <div className="form-group">
        <label>Images (comma separated):</label>
        <input type="text" value={images} onChange={(e) => setImages(e.target.value)} />
      </div>
      <div className="form-group">
        <label>SKU:</label>
        <input type="text" value={sku} onChange={(e) => setSku(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Description:</label>
        <textarea value={description} onChange={(e) => setDescription(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Attributes (key:value, comma separated):</label>
        <input type="text" value={attributes} onChange={(e) => setAttributes(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Fitment Data (comma separated):</label>
        <input type="text" value={fitmentData} onChange={(e) => setFitmentData(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Location:</label>
        <input type="text" value={location} onChange={(e) => setLocation(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Shipment Weight:</label>
        <input type="number" value={shipmentWeight} onChange={(e) => setShipmentWeight(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Shipment Size:</label>
        <input type="text" value={shipmentSize} onChange={(e) => setShipmentSize(e.target.value)} />
      </div>
      <div className="form-group">
        <label>Hazardous:</label>
        <input type="checkbox" checked={shipmentHazardous} onChange={(e) => setShipmentHazardous(e.target.checked)} />
      </div>
      <div className="form-group">
        <label>Fragile:</label>
        <input type="checkbox" checked={shipmentFragile} onChange={(e) => setShipmentFragile(e.target.checked)} />
      </div>
      <div className="form-group">
        <label>Metadata (key:value, comma separated):</label>
        <input type="text" value={metadata} onChange={(e) => setMetadata(e.target.value)} />
      </div>
      <button type="submit">Update Part</button>
    </form>
  );
};

export default UpdatePartForm;
