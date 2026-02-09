import React from 'react';
import { ArrowLeft, Menu } from 'lucide-react';

const Header = ({ title, onBack, showMenu = false }) => {
    return (
        <header style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            padding: '24px 20px',
            background: 'transparent',
            position: 'relative',
            zIndex: 10
        }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                {onBack && (
                    <button
                        onClick={onBack}
                        style={{
                            background: 'none',
                            border: 'none',
                            color: 'var(--text-primary)',
                            cursor: 'pointer',
                            padding: '4px'
                        }}
                    >
                        <ArrowLeft size={24} />
                    </button>
                )}
                <h1 style={{
                    fontSize: '1.25rem',
                    fontWeight: 600,
                    margin: 0,
                    color: 'var(--text-primary)'
                }}>
                    {title}
                </h1>
            </div>

            {showMenu && (
                <button
                    style={{
                        background: 'none',
                        border: 'none',
                        color: 'var(--text-primary)',
                        cursor: 'pointer'
                    }}
                >
                    <Menu size={24} />
                </button>
            )}
        </header>
    );
};

export default Header;
