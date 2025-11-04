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
    const startTime = formData.get('startTime');
    const endTime = formData.get('endTime');
    const date = formData.get('date');
    
    const data = {
        task: formData.get('task').trim(),
        description: formData.get('description').trim(),
        category: formData.get('category'),
        duration: parseInt(formData.get('duration')),
        date: date
    };
    
    // Add start_time and end_time if provided
    if (startTime && endTime) {
        data.start_time = `${date}T${startTime}:00Z`;  // Send full ISO timestamp
        data.end_time = `${date}T${endTime}:00Z`;      // Send full ISO timestamp
    }
    
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
            const errorText = await response.text();
            throw new Error(`Server error (${response.status}): ${errorText}`);
        }
        
        const newEntry = await response.json();
        entries.unshift(newEntry);
        updateTodayStats();
        loadDayEntries(); // Refresh time slots to show the new booking
        
        // Reset form
        event.target.reset();
        document.getElementById('date').value = new Date().toISOString().split('T')[0];
        clearSelectedSlots(); // Clear any selected time slots
        
        showSuccess('Time entry added successfully!');
    } catch (error) {
        console.error('Error creating entry:', error);
        showError(`Failed to create time entry: ${error.message}`);
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

// Time Slots Functionality
let selectedTimeSlots = [];
let isSelecting = false;
let dayEntries = [];

// Initialize time slots when page loads
document.addEventListener('DOMContentLoaded', function() {
    // Setup time slots
    initializeTimeSlots();
    
    // Setup date change listener
    document.getElementById('date').addEventListener('change', handleDateChange);
    
    // Setup time input listeners
    document.getElementById('startTime').addEventListener('change', handleTimeChange);
    document.getElementById('endTime').addEventListener('change', handleTimeChange);
    document.getElementById('duration').addEventListener('change', handleDurationChange);
});

function initializeTimeSlots() {
    generateTimeSlots();
    updateSelectedDateDisplay();
    loadDayEntries();
}

function generateTimeSlots() {
    const container = document.getElementById('timeSlots');
    if (!container) return;
    
    container.innerHTML = '';
    
    // Generate time slots from 6:00 to 22:00 (6 AM to 10 PM) in 30-minute intervals (32 slots)
    for (let hour = 6; hour < 22; hour++) {
        for (let minute = 0; minute < 60; minute += 30) {
            const timeSlot = createTimeSlot(hour, minute);
            container.appendChild(timeSlot);
        }
    }
}

function createTimeSlot(hour, minute) {
    const slot = document.createElement('div');
    slot.className = 'time-slot available';
    
    const timeString = `${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`;
    const slotId = `${hour}-${minute}`;
    slot.dataset.time = timeString;
    slot.dataset.slotId = slotId;
    slot.dataset.hour = hour;
    slot.dataset.minute = minute;
    
    const timeElement = document.createElement('div');
    timeElement.className = 'time-slot-time';
    timeElement.textContent = timeString;
    
    const contentElement = document.createElement('div');
    contentElement.className = 'time-slot-content';
    
    slot.appendChild(timeElement);
    slot.appendChild(contentElement);
    
    // Add click and drag functionality
    slot.addEventListener('mousedown', handleSlotMouseDown);
    slot.addEventListener('mouseenter', handleSlotMouseEnter);
    slot.addEventListener('mouseup', handleSlotMouseUp);
    
    return slot;
}

function handleSlotMouseDown(event) {
    if (event.target.closest('.time-slot').classList.contains('booked')) return;
    
    isSelecting = true;
    selectedTimeSlots = [];
    clearSelectedSlots();
    
    const slot = event.target.closest('.time-slot');
    selectSlot(slot);
    
    event.preventDefault();
}

function handleSlotMouseEnter(event) {
    if (!isSelecting) return;
    if (event.target.closest('.time-slot').classList.contains('booked')) return;
    
    const slot = event.target.closest('.time-slot');
    selectSlot(slot);
}

function handleSlotMouseUp(event) {
    if (!isSelecting) return;
    
    isSelecting = false;
    updateTimeInputsFromSelection();
}

// Global mouse up listener to handle drag outside slots
document.addEventListener('mouseup', function() {
    if (isSelecting) {
        isSelecting = false;
        updateTimeInputsFromSelection();
    }
});

function selectSlot(slot) {
    if (!slot || slot.classList.contains('booked')) return;
    
    const slotId = slot.dataset.slotId;
    if (!selectedTimeSlots.includes(slotId)) {
        selectedTimeSlots.push(slotId);
        slot.classList.add('selected');
        slot.classList.remove('available');
    }
}

function clearSelectedSlots() {
    document.querySelectorAll('.time-slot.selected').forEach(slot => {
        slot.classList.remove('selected');
        if (!slot.classList.contains('booked')) {
            slot.classList.add('available');
        }
    });
    selectedTimeSlots = [];
}

function updateTimeInputsFromSelection() {
    if (selectedTimeSlots.length === 0) return;
    
    // Sort selected slots by time
    const sortedSlots = selectedTimeSlots.sort((a, b) => {
        const [hourA, minuteA] = a.split('-').map(Number);
        const [hourB, minuteB] = b.split('-').map(Number);
        return (hourA * 60 + minuteA) - (hourB * 60 + minuteB);
    });
    
    // Get start time from first slot
    const [startHour, startMinute] = sortedSlots[0].split('-').map(Number);
    const startTime = `${startHour.toString().padStart(2, '0')}:${startMinute.toString().padStart(2, '0')}`;
    
    // Calculate end time (last slot + 30 minutes)
    const [endHour, endMinute] = sortedSlots[sortedSlots.length - 1].split('-').map(Number);
    let endTotalMinutes = endHour * 60 + endMinute + 30;
    const endHourFinal = Math.floor(endTotalMinutes / 60) % 24;
    const endMinuteFinal = endTotalMinutes % 60;
    const endTime = `${endHourFinal.toString().padStart(2, '0')}:${endMinuteFinal.toString().padStart(2, '0')}`;
    
    // Calculate duration
    const startTotalMinutes = startHour * 60 + startMinute;
    const duration = (endTotalMinutes - startTotalMinutes);
    
    // Update form fields
    document.getElementById('startTime').value = startTime;
    document.getElementById('endTime').value = endTime;
    document.getElementById('duration').value = duration;
}

function handleTimeChange() {
    const startTime = document.getElementById('startTime').value;
    const endTime = document.getElementById('endTime').value;
    
    if (startTime && endTime) {
        // Calculate duration
        const start = new Date(`2000-01-01T${startTime}`);
        const end = new Date(`2000-01-01T${endTime}`);
        
        if (end > start) {
            const duration = (end - start) / (1000 * 60); // minutes
            document.getElementById('duration').value = duration;
            
            // Update visual selection
            updateSlotSelectionFromTimes(startTime, endTime);
        }
    }
}

function handleDurationChange() {
    const startTime = document.getElementById('startTime').value;
    const duration = parseInt(document.getElementById('duration').value);
    
    if (startTime && duration) {
        // Calculate end time
        const start = new Date(`2000-01-01T${startTime}`);
        const end = new Date(start.getTime() + duration * 60000);
        const endTime = end.toTimeString().slice(0, 5);
        
        document.getElementById('endTime').value = endTime;
        
        // Update visual selection
        updateSlotSelectionFromTimes(startTime, endTime);
    }
}

function updateSlotSelectionFromTimes(startTime, endTime) {
    clearSelectedSlots();
    
    const [startHour, startMinute] = startTime.split(':').map(Number);
    const [endHour, endMinute] = endTime.split(':').map(Number);
    
    const startTotalMinutes = startHour * 60 + startMinute;
    const endTotalMinutes = endHour * 60 + endMinute;
    
    // Select slots that fall within the time range
    document.querySelectorAll('.time-slot.available').forEach(slot => {
        const hour = parseInt(slot.dataset.hour);
        const minute = parseInt(slot.dataset.minute);
        const slotTotalMinutes = hour * 60 + minute;
        
        if (slotTotalMinutes >= startTotalMinutes && slotTotalMinutes < endTotalMinutes) {
            slot.classList.add('selected');
            slot.classList.remove('available');
            selectedTimeSlots.push(slot.dataset.slotId);
        }
    });
}

function handleDateChange() {
    updateSelectedDateDisplay();
    loadDayEntries();
    clearSelectedSlots();
    
    // Clear time inputs
    document.getElementById('startTime').value = '';
    document.getElementById('endTime').value = '';
    document.getElementById('duration').value = '60';
}

function updateSelectedDateDisplay() {
    const dateInput = document.getElementById('date');
    const displayElement = document.getElementById('selectedDateDisplay');
    
    if (dateInput && displayElement) {
        const selectedDate = new Date(dateInput.value + 'T00:00:00');
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        selectedDate.setHours(0, 0, 0, 0);
        
        if (selectedDate.getTime() === today.getTime()) {
            displayElement.textContent = 'Today';
        } else {
            displayElement.textContent = selectedDate.toLocaleDateString('en-US', {
                weekday: 'long',
                month: 'short',
                day: 'numeric'
            });
        }
    }
}

function loadDayEntries() {
    const selectedDate = document.getElementById('date').value;
    dayEntries = entries.filter(entry => entry.date === selectedDate);
    
    updateTimeSlotsWithBookings();
}

function updateTimeSlotsWithBookings() {
    // First, reset all slots to available
    document.querySelectorAll('.time-slot').forEach(slot => {
        slot.classList.remove('booked', 'selected');
        slot.classList.add('available');
        
        const contentElement = slot.querySelector('.time-slot-content');
        contentElement.innerHTML = '';
    });
    
    // Mark booked slots
    dayEntries.forEach(entry => {
        if (entry.start_time && entry.end_time) {
            const startTime = new Date(entry.start_time);
            const endTime = new Date(entry.end_time);
            
            const startHour = startTime.getHours();
            const startMinute = Math.floor(startTime.getMinutes() / 30) * 30;
            const endHour = endTime.getHours();
            const endMinute = Math.ceil(endTime.getMinutes() / 30) * 30;
            
            // Mark all slots within this time range as booked
            document.querySelectorAll('.time-slot').forEach(slot => {
                const slotHour = parseInt(slot.dataset.hour);
                const slotMinute = parseInt(slot.dataset.minute);
                
                const slotTotalMinutes = slotHour * 60 + slotMinute;
                const startTotalMinutes = startHour * 60 + startMinute;
                const endTotalMinutes = endHour * 60 + endMinute;
                
                if (slotTotalMinutes >= startTotalMinutes && slotTotalMinutes < endTotalMinutes) {
                    slot.classList.remove('available');
                    slot.classList.add('booked');
                    
                    const contentElement = slot.querySelector('.time-slot-content');
                    const taskElement = document.createElement('div');
                    taskElement.className = 'time-slot-task';
                    taskElement.textContent = entry.task;
                    
                    const categoryElement = document.createElement('div');
                    categoryElement.className = 'time-slot-category';
                    categoryElement.textContent = entry.category || 'Other';
                    
                    contentElement.appendChild(taskElement);
                    contentElement.appendChild(categoryElement);
                }
            });
        }
    });
}