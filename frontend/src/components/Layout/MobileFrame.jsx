import React from 'react';
import { motion } from 'framer-motion';

const MobileFrame = ({ children }) => {
    return (
        <div className="mobile-frame-container">
            <motion.div
                className="mobile-frame"
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5, ease: "easeOut" }}
            >
                {children}
            </motion.div>
            <style>{`
        .mobile-frame-container {
          display: flex;
          justify-content: center;
          align-items: center;
          min-height: 100vh;
          width: 100vw;
          background-color: #000;
          padding: 0;
        }

        .mobile-frame {
          width: 100%;
          height: 100%;
          max-width: 480px; /* Max width for mobile simulation */
          background-color: var(--bg-primary);
          position: relative;
          overflow: hidden;
          display: flex;
          flex-direction: column;
          box-shadow: 0 0 0 1px #333;
        }

        @media (min-width: 481px) {
          .mobile-frame {
            height: 90vh;
            border-radius: 40px;
            border: 8px solid #222;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
          }
        }
      `}</style>
        </div>
    );
};

export default MobileFrame;
