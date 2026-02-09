import React from 'react';
import { motion } from 'framer-motion';

const Button = ({
    children,
    variant = 'primary',
    fullWidth = false,
    onClick,
    disabled = false,
    className = ''
}) => {
    const baseStyles = "inline-flex items-center justify-center font-medium transition-all focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed";

    const variants = {
        primary: {
            backgroundColor: 'var(--text-primary)',
            color: 'var(--bg-primary)',
            border: 'none',
            boxShadow: '0 4px 6px -1px rgba(255, 255, 255, 0.1)',
        },
        secondary: {
            backgroundColor: 'var(--bg-tertiary)',
            color: 'var(--text-primary)',
            border: '1px solid var(--bg-tertiary)',
        },
        outline: {
            backgroundColor: 'transparent',
            color: 'var(--text-primary)',
            border: '1px solid var(--text-secondary)',
        },
        ghost: {
            backgroundColor: 'transparent',
            color: 'var(--text-secondary)',
            border: 'none',
        }
    };

    const style = {
        ...variants[variant],
        padding: '12px 24px',
        borderRadius: '12px',
        width: fullWidth ? '100%' : 'auto',
        cursor: disabled ? 'not-allowed' : 'pointer',
        fontSize: '1rem',
        letterSpacing: '0.025em',
    };

    return (
        <motion.button
            whileTap={{ scale: 0.98 }}
            whileHover={{ scale: 1.02 }}
            style={style}
            onClick={onClick}
            disabled={disabled}
            className={`${baseStyles} ${className}`}
        >
            {children}
        </motion.button>
    );
};

export default Button;
