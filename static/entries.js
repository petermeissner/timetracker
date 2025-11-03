// Global variables
let entries = [];
let categories = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Set today's date as default
    const today = new Date().toISOString().split('T')[0];
    
    // Load data
    loadCategories();
    loadEntries();
    
    // Setup event listeners
    setupEventListeners();
    
    // Setup category filter
    document.getElementById('categoryFilter').addEventListener('change', filterEntries);
});

function setupEventListeners() {
    // Edit form submission
    document.getElementById('editForm').addEventListener('submit', handleEditSubmit);
    
    // Modal close
    document.querySelector('.close').addEventListener('click', closeEditModal);
    window.addEventListener('click', function(event) {
        const modal = document.getElementById('editModal');
        if (event.target === modal) {
            closeEditModal();
        }
    });
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

async function loadEntries() {
    try {
        console.log('Loading entries...');
        const response = await fetch('/api/entries');
        if (!response.ok) {
            throw new Error('Failed to load entries');
        }
        
        entries = await response.json() || [];
        console.log('Loaded entries:', entries.length, entries);
        renderEntries();
        updateTotalTime();
    } catch (error) {
        console.error('Error loading entries:', error);
        showError('Failed to load time entries');
    }
}

function updateCategorySelectors() {
    // Update category selectors
    const editCategorySelect = document.getElementById('editCategory');
    const filterSelect = document.getElementById('categoryFilter');
    
    let categoryOptions = '<option value="">Select a category...</option>';
    let filterOptions = '<option value="">All Categories</option>';
    
    categories.forEach(category => {
        categoryOptions += `<option value="${category.name}">${escapeHtml(category.name)}</option>`;
        filterOptions += `<option value="${category.name}">${escapeHtml(category.name)}</option>`;
    });
    
    if (editCategorySelect) editCategorySelect.innerHTML = categoryOptions;
    if (filterSelect) filterSelect.innerHTML = filterOptions;
}

function renderEntries() {
    const entriesList = document.getElementById('entriesList');
    const categoryFilter = document.getElementById('categoryFilter').value;
    
    // Filter entries by category if a filter is selected
    let filteredEntries = entries;
    if (categoryFilter) {
        filteredEntries = entries.filter(entry => entry.category === categoryFilter);
    }
    
    if (filteredEntries.length === 0) {
        const message = categoryFilter ? 
            `No entries found for category: ${getCategoryInfo(categoryFilter).name}` :
            'No time entries yet';
        const subMessage = categoryFilter ? 
            'Try selecting a different category or clear the filter.' :
            'Go back to the main page to add your first time entry!';
            
        entriesList.innerHTML = `
            <div class="empty-state">
                <h3>${message}</h3>
                <p>${subMessage}</p>
            </div>
        `;
        return;
    }
    
    // Group entries by date
    const groupedEntries = groupEntriesByDate(filteredEntries);
    
    let html = '';
    for (const [date, dateEntries] of Object.entries(groupedEntries)) {
        const totalMinutes = dateEntries.reduce((sum, entry) => sum + entry.duration, 0);
        const totalHours = Math.floor(totalMinutes / 60);
        const remainingMinutes = totalMinutes % 60;
        
        html += `
            <div class="date-group">
                <h3 class="date-header">
                    ${formatDate(date)} 
                    <span class="date-total">(${totalHours}h ${remainingMinutes}m)</span>
                </h3>
                <div class="date-entries">
        `;
        
        dateEntries.forEach(entry => {
            html += renderEntry(entry);
        });
        
        html += `
                </div>
            </div>
        `;
    }
    
    entriesList.innerHTML = html;
}

function renderEntry(entry) {
    const hours = Math.floor(entry.duration / 60);
    const minutes = entry.duration % 60;
    const durationText = hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
    
    // Get category display name and color
    const categoryInfo = getCategoryInfo(entry.category || 'other');
    
    return `
        <div class="entry-item" data-id="${entry.id}">
            <div class="entry-header">
                <div class="entry-task">${escapeHtml(entry.task)}</div>
                <div class="entry-badges">
                    <div class="category-badge" style="background-color: ${categoryInfo.color}; color: white;">${categoryInfo.name}</div>
                    <div class="duration-badge">${durationText}</div>
                </div>
            </div>
            
            <div class="entry-meta">
                <span>📅 ${formatDate(entry.date)}</span>
                ${entry.start_time ? `<span>⏰ ${formatTime(entry.start_time)} - ${formatTime(entry.end_time)}</span>` : ''}
            </div>
            
            ${entry.description ? `<div class="entry-description">${escapeHtml(entry.description)}</div>` : ''}
            
            <div class="entry-actions">
                <button class="btn btn-secondary btn-small" onclick="editEntry(${entry.id})">Edit</button>
                <button class="btn btn-danger btn-small" onclick="deleteEntry(${entry.id})">Delete</button>
            </div>
        </div>
    `;
}

function groupEntriesByDate(entries) {
    return entries.reduce((groups, entry) => {
        const date = entry.date;
        if (!groups[date]) {
            groups[date] = [];
        }
        groups[date].push(entry);
        return groups;
    }, {});
}

function updateTotalTime() {
    const today = new Date().toISOString().split('T')[0];
    const categoryFilter = document.getElementById('categoryFilter').value;
    
    let todayEntries = entries.filter(entry => entry.date === today);
    if (categoryFilter) {
        todayEntries = todayEntries.filter(entry => entry.category === categoryFilter);
    }
    
    const totalMinutes = todayEntries.reduce((sum, entry) => sum + entry.duration, 0);
    
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;
    const totalText = hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
    
    const filterText = categoryFilter ? ` (${getCategoryInfo(categoryFilter).name})` : '';
    document.getElementById('totalTime').textContent = totalText + filterText;
}

function filterEntries() {
    renderEntries();
    updateTotalTime();
}

async function editEntry(id) {
    const entry = entries.find(e => e.id === id);
    if (!entry) return;
    
    // Populate edit form
    document.getElementById('editId').value = entry.id;
    document.getElementById('editTask').value = entry.task;
    document.getElementById('editCategory').value = entry.category || 'other';
    document.getElementById('editDescription').value = entry.description || '';
    document.getElementById('editDuration').value = entry.duration;
    document.getElementById('editDate').value = entry.date;
    
    // Show modal
    document.getElementById('editModal').style.display = 'block';
}

async function handleEditSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const id = parseInt(formData.get('id') || document.getElementById('editId').value);
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
        const response = await fetch(`/api/entries/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error('Failed to update entry');
        }
        
        const updatedEntry = await response.json();
        const index = entries.findIndex(e => e.id === id);
        if (index !== -1) {
            entries[index] = updatedEntry;
        }
        
        renderEntries();
        updateTotalTime();
        closeEditModal();
        
        showSuccess('Time entry updated successfully!');
    } catch (error) {
        console.error('Error updating entry:', error);
        showError(`Failed to update time entry: ${error.message}`);
    }
}

async function deleteEntry(id) {
    console.log('deleteEntry called with id:', id);
    
    if (!confirm('Are you sure you want to delete this time entry?')) {
        console.log('Delete cancelled by user');
        return;
    }
    
    console.log('Attempting to delete entry with id:', id);
    
    try {
        const response = await fetch(`/api/entries/${id}`, {
            method: 'DELETE'
        });
        
        console.log('Delete response status:', response.status);
        console.log('Delete response ok:', response.ok);
        
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Delete failed with response:', errorText);
            throw new Error(`Failed to delete entry: ${response.status} ${errorText}`);
        }
        
        console.log('Delete successful, updating UI');
        entries = entries.filter(e => e.id !== id);
        renderEntries();
        updateTotalTime();
        
        showSuccess('Time entry deleted successfully!');
    } catch (error) {
        console.error('Error deleting entry:', error);
        showError(`Failed to delete time entry: ${error.message}`);
    }
}

function closeEditModal() {
    document.getElementById('editModal').style.display = 'none';
    document.getElementById('editForm').reset();
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