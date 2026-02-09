import React from 'react';
import { motion } from 'framer-motion';
import { ChevronRight, Clock, AlertCircle } from 'lucide-react';
import Header from '../components/UI/Header';
import Card from '../components/UI/Card';

const Dashboard = ({ onSelectProject }) => {
    const [projects, setProjects] = React.useState([]);
    const [loading, setLoading] = React.useState(true);
    const [error, setError] = React.useState(null);

    React.useEffect(() => {
        fetch('http://localhost:8080/api/projects')
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch projects');
                return res.json();
            })
            .then(data => {
                setProjects(data);
                setLoading(false);
            })
            .catch(err => {
                console.error("API Error:", err);
                setError(err.message);
                setLoading(false);
                // Fallback to mock data for demo purposes if backend is down
                /* 
                setProjects([
                     { id: 1, name: 'WebApp-Landing (Offline)', status: 'Backend unreachable', state: 'error' },
                ]); 
                */
            });
    }, []);

    if (loading) return <div style={{ padding: 20, color: '#fff' }}>Loading projects...</div>;
    if (error) return <div style={{ padding: 20, color: 'var(--status-error)' }}>Error: {error}. Is the backend running?</div>;

    return (
        <>
            <Header title="Projects" onBack={() => { }} showMenu />

            <div style={{ padding: '0 20px 80px', display: 'flex', flexDirection: 'column', gap: '16px' }}>
                {projects.map((project, index) => (
                    <motion.div
                        key={project.id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                    >
                        <Card onClick={() => onSelectProject(project)}>
                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <div>
                                    <h3 style={{ fontSize: '1.25rem', fontWeight: 600, marginBottom: '4px' }}>{project.name}</h3>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                                        {project.state === 'error' && <AlertCircle size={14} color="var(--status-error)" />}
                                        {project.state === 'success' && <div style={{ width: 8, height: 8, borderRadius: '50%', background: 'var(--status-success)' }} />}
                                        <span style={{
                                            fontSize: '0.875rem',
                                            color: project.state === 'error' ? 'var(--status-error)' : 'var(--text-secondary)'
                                        }}>
                                            {project.status}
                                        </span>
                                    </div>
                                </div>
                                <ChevronRight size={20} color="var(--text-muted)" />
                            </div>
                        </Card>
                    </motion.div>
                ))}
            </div>
        </>
    );
};

export default Dashboard;
