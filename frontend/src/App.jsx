import React, { useState } from 'react';
import MobileFrame from './components/Layout/MobileFrame';
import Login from './views/Login';
import Dashboard from './views/Dashboard';
import ProjectControl from './views/ProjectControl';
import BuildProgress from './views/BuildProgress';
import Preview from './views/Preview';

function App() {
  const [currentView, setCurrentView] = useState('login');
  const [selectedProject, setSelectedProject] = useState(null);

  const handleLogin = () => {
    setCurrentView('dashboard');
  };

  const handleSelectProject = (project) => {
    setSelectedProject(project);
    setCurrentView('control');
  };

  const handleStartBuild = () => {
    setCurrentView('build');
  };

  const handleApproveDeploy = () => {
    setCurrentView('preview');
  };

  const handleBack = () => {
    if (currentView === 'control') setCurrentView('dashboard');
    if (currentView === 'build') setCurrentView('control');
    if (currentView === 'preview') setCurrentView('build');
  };

  return (
    <MobileFrame>
      {currentView === 'login' && <Login onLogin={handleLogin} />}

      {currentView === 'dashboard' && (
        <Dashboard onSelectProject={handleSelectProject} />
      )}

      {currentView === 'control' && (
        <ProjectControl
          project={selectedProject}
          onBack={handleBack}
          onStartBuild={handleStartBuild}
        />
      )}

      {currentView === 'build' && (
        <BuildProgress
          project={selectedProject}
          onBack={handleBack}
          onApprove={handleApproveDeploy}
        />
      )}

      {currentView === 'preview' && (
        <Preview onBack={handleBack} />
      )}
    </MobileFrame>
  );
}

export default App;
