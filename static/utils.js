/**
 * Shared utility functions for the timesheet application
 * Centralized location for common functionality across all pages
 */

const Utils = {
    /**
     * Date and Time Formatting
     */
    
    formatDate(dateString) {
        const date = new Date(dateString + 'T00:00:00');
        return date.toLocaleDateString('en-US', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    },

    formatDateShort(dateString) {
        const date = new Date(dateString + 'T00:00:00');
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);
        
        // Check if it's today
        if (date.toDateString() === today.toDateString()) {
            return 'Today';
        }
        
        // Check if it's yesterday
        if (date.toDateString() === yesterday.toDateString()) {
            return 'Yesterday';
        }
        
        // For other dates, show short format
        return date.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: date.getFullYear() !== today.getFullYear() ? 'numeric' : undefined
        });
    },

    formatTime(timeString) {
        if (!timeString) return '';
        // Parse times without timezone conversion (treat as local time)
        const timeStr = timeString.replace('Z', '').replace('T', ' ');
        const date = new Date(timeStr);
        return date.toLocaleTimeString('en-US', {
            hour: '2-digit',
            minute: '2-digit',
            hour12: false  // Use 24-hour military time format
        });
    },

    /**
     * HTML Safety and Text Processing
     */
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    },

    /**
     * Notification System
     */
    
    showSuccess(message) {
        const notification = document.createElement('div');
        notification.className = 'notification success';
        notification.textContent = message;
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: #48bb78;
            color: white;
            padding: 15px 20px;
            border-radius: 6px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            z-index: 1001;
            animation: slideIn 0.3s ease-out;
        `;
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.remove();
        }, 3000);
    },

    showError(message) {
        const notification = document.createElement('div');
        notification.className = 'notification error';
        notification.textContent = message;
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: #e53e3e;
            color: white;
            padding: 15px 20px;
            border-radius: 6px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            z-index: 1001;
            animation: slideIn 0.3s ease-out;
        `;
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.remove();
        }, 5000);
    },

    /**
     * Category Helpers
     */
    
    getCategoryInfo(categoryName, categories = []) {
        const category = categories.find(c => c.name === categoryName);
        if (category) {
            return {
                name: category.name.split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' '),
                class: `category-${category.name.replace(/\s+/g, '-')}`,
                color: category.color
            };
        }
        
        // Fallback for unknown categories
        return {
            name: categoryName ? categoryName.split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ') : 'Other',
            class: 'category-other',
            color: '#718096'
        };
    }
};

// Initialize notification animations CSS if not already present
if (!document.querySelector('#utils-styles')) {
    const style = document.createElement('style');
    style.id = 'utils-styles';
    style.textContent = `
        @keyframes slideIn {
            from {
                transform: translateX(100%);
                opacity: 0;
            }
            to {
                transform: translateX(0);
                opacity: 1;
            }
        }
    `;
    document.head.appendChild(style);
}

// Export for use in other modules
window.Utils = Utils;