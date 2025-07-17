import React from 'react';
import { useParams } from 'react-router-dom';

const ProjectView: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  return (
    <div style={{ padding: '2rem' }}>
      <h1>Project View</h1>
      <p>Project ID: {id}</p>
      <p>This is a placeholder for the project view component.</p>
    </div>
  );
};

export default ProjectView;