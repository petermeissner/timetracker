// Global variables
let entries = [];
let editingId = null;
let categories = [];
let predefinedTasks = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Set today's date as default
    const today = new Date().toISOString().split('T')[0];
    document.getElementById('date').value = today;
    
    // Load data
    loadCategories();
    loadPredefinedTasks();
    loadEntries();
    
    // Setup event listeners
    setupEventListeners();
    
    // Setup predefined task selector
    document.getElementById('predefinedTask').addEventListener('change', handlePredefinedTaskSelect);
});

function setupEventListeners() {
    // Form submission
    document.getElementById('timeEntryForm').addEventListener('submit', handleFormSubmit);
}

async function loadCategories() {
    try {
        const response = await fetch('/api/categories');
        if (!response.ok) {
            throw new Error('Failed to load categories');
        }
        
        categories = await response.json() || [];
        updateCategorySelectors();
    } catch (error) {
        console.error('Error loading categories:', error);
        // Continue with default categories if API fails
        categories = [
            { id: 1, name: 'project work', color: '#48bb78' },
            { id: 2, name: 'project support', color: '#ed8936' },
            { id: 3, name: 'other', color: '#718096' }
        ];
        updateCategorySelectors();
    }
}

async function loadPredefinedTasks() {
    try {
        const response = await fetch('/api/tasks');
        if (!response.ok) {
            throw new Error('Failed to load tasks');
        }
        
        predefinedTasks = await response.json() || [];
        updatePredefinedTaskSelector();
    } catch (error) {
        console.error('Error loading predefined tasks:', error);
        predefinedTasks = [];
        updatePredefinedTaskSelector();
    }
}

async function loadEntries() {
    try {
        console.log('Loading entries...');
        const response = await fetch('/api/entries');
        if (!response.ok) {
            throw new Error('Failed to load entries');
        }
        
        entries = await response.json() || [];
        console.log('Loaded entries:', entries.length, entries);
        updateTodayStats();
    } catch (error) {
        console.error('Error loading entries:', error);
        showError('Failed to load time entries');
    }
}

function updateCategorySelectors() {
    // Update main category selector
    const categorySelect = document.getElementById('category');
    const editCategorySelect = document.getElementById('editCategory');
    const filterSelect = document.getElementById('categoryFilter');
    
    let categoryOptions = '<option value="">Select a category...</option>';
    
    categories.forEach(category => {
        categoryOptions += `<option value="${category.name}">${escapeHtml(category.name)}</option>`;
    });
    
    if (categorySelect) categorySelect.innerHTML = categoryOptions;
}

function updatePredefinedTaskSelector() {
    const taskSelect = document.getElementById('predefinedTask');
    
    let options = '<option value="">Select predefined task...</option>';
    
    // Group tasks by category
    const tasksByCategory = {};
    predefinedTasks.forEach(task => {
        const categoryName = categories.find(c => c.id === task.category_id)?.name || 'Uncategorized';
        if (!tasksByCategory[categoryName]) {
            tasksByCategory[categoryName] = [];
        }
        tasksByCategory[categoryName].push(task);
    });
    
    // Add options grouped by category
    Object.keys(tasksByCategory).sort().forEach(categoryName => {
        if (tasksByCategory[categoryName].length > 0) {
            options += `<optgroup label="${escapeHtml(categoryName)}">`;
            tasksByCategory[categoryName].forEach(task => {
                options += `<option value="${task.id}" data-category="${task.category_id || ''}" data-description="${escapeHtml(task.description || '')}">${escapeHtml(task.name)}</option>`;
            });
            options += '</optgroup>';
        }
    });
    
    if (taskSelect) taskSelect.innerHTML = options;
}

function handlePredefinedTaskSelect(event) {
    const selectedOption = event.target.options[event.target.selectedIndex];
    if (!selectedOption.value) return;
    
    const task = predefinedTasks.find(t => t.id === parseInt(selectedOption.value));
    if (!task) return;
    
    // Fill in the form fields
    document.getElementById('task').value = task.name;
    
    if (task.category_id) {
        const category = categories.find(c => c.id === task.category_id);
        if (category) {
            document.getElementById('category').value = category.name;
        }
    }
    
    if (task.description) {
        document.getElementById('description').value = task.description;
    }
    
    // Reset the selector
    event.target.selectedIndex = 0;
}



function updateTodayStats() {
    const today = new Date();
    
    // Format time strings
    const formatTime = (minutes) => {
        if (minutes === 0) return '0m';
        const hours = Math.floor(minutes / 60);
        const mins = minutes % 60;
        if (hours > 0 && mins > 0) {
            return `${hours}h ${mins}m`;
        } else if (hours > 0) {
            return `${hours}h`;
        } else {
            return `${mins}m`;
        }
    };
    
    // Generate the last 7 days (including today)
    const days = [];
    
    for (let i = 6; i >= 0; i--) {
        const date = new Date(today);
        date.setDate(date.getDate() - i);
        const dateStr = date.toISOString().split('T')[0];
        
        // Calculate total for this day
        const dayEntries = entries.filter(entry => entry.date === dateStr);
        const dayMinutes = dayEntries.reduce((sum, entry) => sum + entry.duration, 0);
        
        // Format weekday and date
        const weekdayName = date.toLocaleDateString('en-US', { weekday: 'long' });
        const shortDate = date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
        
        // Determine if this is today
        const isToday = dateStr === today.toISOString().split('T')[0];
        
        days.push({
            date: dateStr,
            weekday: weekdayName,
            shortDate: shortDate,
            minutes: dayMinutes,
            isToday: isToday
        });
    }
    
    // Generate HTML for the breakdown
    const breakdownContainer = document.getElementById('sevenDayBreakdown');
    if (breakdownContainer) {
        let html = '';
        
        days.forEach(day => {
            const date = new Date(day.date);
            const dayOfWeek = date.getDay(); // 0 = Sunday, 6 = Saturday
            const isWeekend = dayOfWeek === 0 || dayOfWeek === 6; // Sunday or Saturday
            const isLowTime = day.minutes < 300 && day.minutes > 0 && !isWeekend; // Less than 5 hours, only for weekdays
            
            // Build CSS classes
            let dayClasses = ['day-item'];
            
            if (day.isToday) {
                dayClasses.push('today');
            } else if (isWeekend) {
                dayClasses.push('weekend');
            }
            
            if (isLowTime) {
                dayClasses.push('low-time');
            }
            
            // Use abbreviated weekday names for compact display
            const shortWeekday = day.weekday.substring(0, 3); // Mon, Tue, Wed, etc.
            
            html += `
                <div class="${dayClasses.join(' ')}">
                    <div class="day-info">
                        <div class="day-name">${shortWeekday}</div>
                        <div class="day-date">${day.shortDate}</div>
                    </div>
                    <div class="day-time">${formatTime(day.minutes)}</div>
                </div>
            `;
        });
        
        breakdownContainer.innerHTML = html;
    }
}

async function handleFormSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const data = {
        task: formData.get('task').trim(),
        description: formData.get('description').trim(),
        category: formData.get('category'),
        duration: parseInt(formData.get('duration')),
        date: formData.get('date')
    };
    
    if (!data.task || !data.category || !data.duration || !data.date) {
        showError('Please fill in all required fields');
        return;
    }
    
    try {
        const response = await fetch('/api/entries', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error('Failed to create entry');
        }
        
        const newEntry = await response.json();
        entries.unshift(newEntry);
        updateTodayStats();
        
        // Reset form
        event.target.reset();
        document.getElementById('date').value = new Date().toISOString().split('T')[0];
        
        showSuccess('Time entry added successfully!');
    } catch (error) {
        console.error('Error creating entry:', error);
        showError('Failed to create time entry');
    }
}



// Utility functions
function getCategoryInfo(categoryName) {
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

function formatDate(dateString) {
    const date = new Date(dateString + 'T00:00:00');
    return date.toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

function formatTime(timeString) {
    if (!timeString) return '';
    const date = new Date(timeString);
    return date.toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit'
    });
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showSuccess(message) {
    // Simple success notification
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
}

function showError(message) {
    // Simple error notification
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
}

// Add CSS for notifications
const style = document.createElement('style');
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
    
    .date-group {
        margin-bottom: 30px;
    }
    
    .date-header {
        color: #4a5568;
        margin-bottom: 15px;
        padding-bottom: 8px;
        border-bottom: 1px solid #e2e8f0;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    
    .date-total {
        font-size: 0.9rem;
        color: #718096;
        font-weight: normal;
    }
    
    .date-entries {
        display: grid;
        gap: 10px;
    }
`;
document.head.appendChild(style);