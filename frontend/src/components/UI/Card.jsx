import React from 'react';
import { motion } from 'framer-motion';

const Card = ({ children, className = '', onClick }) => {
    const styles = {
        background: 'rgba(23, 23, 23, 0.6)',
        backdropFilter: 'blur(12px)',
        border: '1px solid rgba(255, 255, 255, 0.08)',
        borderRadius: 'var(--radius-lg)',
        padding: '20px',
        boxShadow: '0 4px 20px -5px rgba(0, 0, 0, 0.3)',
        color: 'var(--text-primary)'
    };

    return (
        <motion.div
            style={styles}
            className={className}
            onClick={onClick}
            whileHover={onClick ? { scale: 1.01, backgroundColor: 'rgba(30, 30, 30, 0.8)' } : {}}
            transition={{ type: 'spring', stiffness: 300, damping: 20 }}
        >
            {children}
        </motion.div>
    );
};

export default Card;
