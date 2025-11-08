// Global variables
let entries = [];
let categories = [];
let currentSort = { column: 'date', direction: 'desc' };

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Set today's date as default
    const today = new Date().toISOString().split('T')[0];
    
    // Load data
    loadCategories();
    loadEntries();
    
    // Setup event listeners
    setupEventListeners();
    
    // Initialize sort indicators
    updateSortIndicators();
    
    // Setup filters
    document.getElementById('categoryFilter').addEventListener('change', filterEntries);
    document.getElementById('dateFromFilter').addEventListener('change', filterEntries);
    document.getElementById('dateToFilter').addEventListener('change', filterEntries);
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
        categories = await API.categories.getAll() || [];
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
        entries = await API.entries.getAll() || [];
        console.log('Loaded entries:', entries.length, entries);
        renderEntries();
        updateTotalTime();
    } catch (error) {
        console.error('Error loading entries:', error);
        Utils.showError('Failed to load time entries');
    }
}

function updateCategorySelectors() {
    // Update category selectors
    const editCategorySelect = document.getElementById('editCategory');
    const filterSelect = document.getElementById('categoryFilter');
    
    let categoryOptions = '<option value="">Select a category...</option>';
    let filterOptions = '<option value="">All Categories</option>';
    
    categories.forEach(category => {
        categoryOptions += `<option value="${category.name}">${Utils.escapeHtml(category.name)}</option>`;
        filterOptions += `<option value="${category.name}">${Utils.escapeHtml(category.name)}</option>`;
    });
    
    if (editCategorySelect) editCategorySelect.innerHTML = categoryOptions;
    if (filterSelect) filterSelect.innerHTML = filterOptions;
}

function renderEntries() {
    const entriesTable = document.getElementById('entriesTable');
    const categoryFilter = document.getElementById('categoryFilter').value;
    const dateFromFilter = document.getElementById('dateFromFilter').value;
    const dateToFilter = document.getElementById('dateToFilter').value;
    
    // Filter entries by category and date range
    let filteredEntries = entries;
    
    // Category filter
    if (categoryFilter) {
        filteredEntries = filteredEntries.filter(entry => entry.category === categoryFilter);
    }
    
    // Date range filter
    if (dateFromFilter || dateToFilter) {
        filteredEntries = filteredEntries.filter(entry => {
            const entryDateStr = Utils.getEntryDate(entry);
            if (!entryDateStr) return false;
            const entryDate = new Date(entryDateStr);
            let matchesFrom = true;
            let matchesTo = true;
            
            if (dateFromFilter) {
                const fromDate = new Date(dateFromFilter);
                matchesFrom = entryDate >= fromDate;
            }
            
            if (dateToFilter) {
                const toDate = new Date(dateToFilter);
                matchesTo = entryDate <= toDate;
            }
            
            return matchesFrom && matchesTo;
        });
    }
    
    if (filteredEntries.length === 0) {
        let message = 'No entries found';
        let subMessage = 'Try adjusting your filters or clear them to see all entries.';
        
        if (!categoryFilter && !dateFromFilter && !dateToFilter) {
            message = 'No time entries yet';
            subMessage = 'Go back to the main page to add your first time entry!';
        }
            
        entriesTable.innerHTML = `
            <div class="empty-state">
                <h3>${message}</h3>
                <p>${subMessage}</p>
            </div>
        `;
        return;
    }
    
    // Apply current sort if any, otherwise default to date descending
    const sortedEntries = applySortToEntries(filteredEntries);
    
    // Create table
    let html = `
        <table class="entries-data-table">
            <thead>
                <tr>
                    <th class="sortable" onclick="sortTable('date')" data-sort="date">
                        Date <span class="sort-indicator" id="sort-date"></span>
                    </th>
                    <th class="sortable" onclick="sortTable('task')" data-sort="task">
                        Task <span class="sort-indicator" id="sort-task"></span>
                    </th>
                    <th class="sortable" onclick="sortTable('category')" data-sort="category">
                        Category <span class="sort-indicator" id="sort-category"></span>
                    </th>
                    <th class="sortable" onclick="sortTable('duration')" data-sort="duration">
                        Duration <span class="sort-indicator" id="sort-duration"></span>
                    </th>
                    <th class="sortable" onclick="sortTable('time')" data-sort="time">
                        Time <span class="sort-indicator" id="sort-time"></span>
                    </th>
                    <th class="sortable" onclick="sortTable('description')" data-sort="description">
                        Description <span class="sort-indicator" id="sort-description"></span>
                    </th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
    `;
    
    sortedEntries.forEach(entry => {
        html += renderTableRow(entry);
    });
    
    html += `
            </tbody>
        </table>
    `;
    
    entriesTable.innerHTML = html;
}

function renderTableRow(entry) {
    const duration = Utils.getEntryDuration(entry);
    const hours = Math.floor(duration / 60);
    const minutes = duration % 60;
    const durationText = hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
    
    // Get category display name and color
    const categoryInfo = Utils.getCategoryInfo(entry.category || 'other', categories);
    
    // Format time range or show "Manual entry" if no start/end time
    const timeRange = entry.start_time && entry.end_time ? 
        `${Utils.formatTime(entry.start_time)} - ${Utils.formatTime(entry.end_time)}` : 
        'Manual';
    
    return `
        <tr data-id="${entry.id}">
            <td class="date-cell">${Utils.formatDateShort(Utils.getEntryDate(entry))}</td>
            <td class="task-cell">${Utils.escapeHtml(entry.task)}</td>
            <td class="category-cell">
                <span class="category-badge" style="background-color: ${categoryInfo.color}; color: white;">
                    ${categoryInfo.name}
                </span>
            </td>
            <td class="duration-cell">${durationText}</td>
            <td class="time-cell">${timeRange}</td>
            <td class="description-cell">${entry.description ? Utils.escapeHtml(entry.description) : '-'}</td>
            <td class="actions-cell">
                <button class="btn btn-secondary btn-small" onclick="editEntry(${entry.id})" title="Edit entry">✏️</button>
                <button class="btn btn-danger btn-small" onclick="deleteEntry(${entry.id})" title="Delete entry">🗑️</button>
            </td>
        </tr>
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
    
    let todayEntries = entries.filter(entry => Utils.getEntryDate(entry) === today);
    if (categoryFilter) {
        todayEntries = todayEntries.filter(entry => entry.category === categoryFilter);
    }
    
    const totalMinutes = todayEntries.reduce((sum, entry) => sum + Utils.getEntryDuration(entry), 0);
    
    const hours = Math.floor(totalMinutes / 60);
    const minutes = totalMinutes % 60;
    const totalText = hours > 0 ? `${hours}h ${minutes}m` : `${minutes}m`;
    
    const filterText = categoryFilter ? ` (${Utils.getCategoryInfo(categoryFilter, categories).name})` : '';
    document.getElementById('totalTime').textContent = totalText + filterText;
}

function filterEntries() {
    renderEntries();
    updateTotalTime();
}

function sortTable(column) {
    // Toggle direction if clicking the same column, otherwise set to ascending
    if (currentSort.column === column) {
        currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentSort.column = column;
        currentSort.direction = 'asc';
    }
    
    renderEntries();
    updateSortIndicators();
}

function applySortToEntries(entries) {
    return entries.sort((a, b) => {
        let aVal, bVal;
        
        switch (currentSort.column) {
            case 'date':
                aVal = new Date(Utils.getEntryDate(a) || '1970-01-01');
                bVal = new Date(Utils.getEntryDate(b) || '1970-01-01');
                break;
            case 'task':
                aVal = a.task.toLowerCase();
                bVal = b.task.toLowerCase();
                break;
            case 'category':
                aVal = (a.category || 'other').toLowerCase();
                bVal = (b.category || 'other').toLowerCase();
                break;
            case 'duration':
                aVal = Utils.getEntryDuration(a);
                bVal = Utils.getEntryDuration(b);
                break;
            case 'time':
                // Sort by start_time if available, otherwise put manual entries at end
                if (a.start_time && b.start_time) {
                    const aTimeStr = a.start_time.replace('Z', '').replace('T', ' ');
                    const bTimeStr = b.start_time.replace('Z', '').replace('T', ' ');
                    aVal = new Date(aTimeStr);
                    bVal = new Date(bTimeStr);
                } else if (a.start_time && !b.start_time) {
                    return currentSort.direction === 'asc' ? -1 : 1;
                } else if (!a.start_time && b.start_time) {
                    return currentSort.direction === 'asc' ? 1 : -1;
                } else {
                    aVal = a.id;
                    bVal = b.id;
                }
                break;
            case 'description':
                aVal = (a.description || '').toLowerCase();
                bVal = (b.description || '').toLowerCase();
                break;
            default:
                aVal = a.id;
                bVal = b.id;
        }
        
        let result;
        if (typeof aVal === 'string') {
            result = aVal.localeCompare(bVal);
        } else if (aVal instanceof Date && bVal instanceof Date) {
            result = aVal - bVal;
        } else {
            result = aVal - bVal;
        }
        
        return currentSort.direction === 'asc' ? result : -result;
    });
}

function updateSortIndicators() {
    // Clear all indicators
    const indicators = document.querySelectorAll('.sort-indicator');
    indicators.forEach(indicator => {
        indicator.textContent = '';
        indicator.parentElement.classList.remove('sorted-asc', 'sorted-desc');
    });
    
    // Set current sort indicator
    const currentIndicator = document.getElementById(`sort-${currentSort.column}`);
    if (currentIndicator) {
        currentIndicator.textContent = currentSort.direction === 'asc' ? ' ↑' : ' ↓';
        currentIndicator.parentElement.classList.add(`sorted-${currentSort.direction}`);
    }
}

function clearDateFilter() {
    document.getElementById('dateFromFilter').value = '';
    document.getElementById('dateToFilter').value = '';
    filterEntries();
}

async function editEntry(id) {
    const entry = entries.find(e => e.id === id);
    if (!entry) return;
    
    // Populate edit form
    document.getElementById('editId').value = entry.id;
    document.getElementById('editTask').value = entry.task;
    document.getElementById('editCategory').value = entry.category || 'other';
    document.getElementById('editDescription').value = entry.description || '';
    document.getElementById('editDuration').value = Utils.getEntryDuration(entry);
    document.getElementById('editDate').value = Utils.getEntryDate(entry);
    
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
        Utils.showError('Please fill in all required fields');
        return;
    }
    
    try {
        const updatedEntry = await API.entries.update(id, data);
        const index = entries.findIndex(e => e.id === id);
        if (index !== -1) {
            entries[index] = updatedEntry;
        }
        
        renderEntries();
        updateTotalTime();
        closeEditModal();
        
        Utils.showSuccess('Time entry updated successfully!');
    } catch (error) {
        console.error('Error updating entry:', error);
        Utils.showError(`Failed to update time entry: ${error.message}`);
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
        await API.entries.delete(id);
        
        console.log('Delete successful, updating UI');
        entries = entries.filter(e => e.id !== id);
        renderEntries();
        updateTotalTime();
        
        Utils.showSuccess('Time entry deleted successfully!');
    } catch (error) {
        console.error('Error deleting entry:', error);
        Utils.showError(`Failed to delete time entry: ${error.message}`);
    }
}

function closeEditModal() {
    document.getElementById('editModal').style.display = 'none';
    document.getElementById('editForm').reset();
}

// Excel Export Functionality
function exportToExcel() {
    try {
        // Get filtered entries (same logic as renderEntries)
        const categoryFilter = document.getElementById('categoryFilter').value;
        const dateFromFilter = document.getElementById('dateFromFilter').value;
        const dateToFilter = document.getElementById('dateToFilter').value;
        
        let filteredEntries = entries;
        
        // Apply category filter
        if (categoryFilter) {
            filteredEntries = filteredEntries.filter(entry => entry.category === categoryFilter);
        }
        
        // Apply date range filter
        if (dateFromFilter || dateToFilter) {
            filteredEntries = filteredEntries.filter(entry => {
                const entryDateStr = Utils.getEntryDate(entry);
                if (!entryDateStr) return false;
                const entryDate = new Date(entryDateStr);
                let matchesFrom = true;
                let matchesTo = true;
                
                if (dateFromFilter) {
                    const fromDate = new Date(dateFromFilter);
                    matchesFrom = entryDate >= fromDate;
                }
                
                if (dateToFilter) {
                    const toDate = new Date(dateToFilter);
                    matchesTo = entryDate <= toDate;
                }
                
                return matchesFrom && matchesTo;
            });
        }
        
        if (filteredEntries.length === 0) {
            Utils.showError('No entries to export with current filters');
            return;
        }
        
        // Sort entries by date (newest first)
        const sortedEntries = filteredEntries.sort((a, b) => {
            return new Date(Utils.getEntryDate(b) || '1970-01-01') - new Date(Utils.getEntryDate(a) || '1970-01-01');
        });
        
        // Prepare data for Excel
        const excelData = [];
        
        // Add header row
        excelData.push([
            'Date',
            'Weekday',
            'Task',
            'Category',
            'Description',
            'Duration (Minutes)',
            'Duration (Hours)',
            'Start Time',
            'End Time'
        ]);
        
        // Add data rows
        sortedEntries.forEach(entry => {
            const entryDateStr = Utils.getEntryDate(entry);
            const entryDate = new Date(entryDateStr + 'T00:00:00');
            const weekday = entryDate.toLocaleDateString('en-US', { weekday: 'long' });
            const duration = Utils.getEntryDuration(entry);
            const hours = (duration / 60).toFixed(2);
            const startTime = entry.start_time ? Utils.formatTime(entry.start_time) : 'Manual';
            const endTime = entry.end_time ? Utils.formatTime(entry.end_time) : 'Manual';
            
            excelData.push([
                entryDateStr,
                weekday,
                entry.task,
                entry.category || 'Other',
                entry.description || '',
                duration,
                hours,
                startTime,
                endTime
            ]);
        });
        
        // Calculate totals
        const totalMinutes = sortedEntries.reduce((sum, entry) => sum + Utils.getEntryDuration(entry), 0);
        const totalHours = (totalMinutes / 60).toFixed(2);
        
        // Add summary rows
        excelData.push([]); // Empty row
        excelData.push(['SUMMARY']);
        excelData.push(['Total Entries:', sortedEntries.length]);
        excelData.push(['Total Minutes:', totalMinutes]);
        excelData.push(['Total Hours:', totalHours]);
        
        // Add filter information
        if (categoryFilter || dateFromFilter || dateToFilter) {
            excelData.push([]); // Empty row
            excelData.push(['FILTERS APPLIED']);
            if (categoryFilter) {
                excelData.push(['Category:', categoryFilter]);
            }
            if (dateFromFilter) {
                excelData.push(['From Date:', dateFromFilter]);
            }
            if (dateToFilter) {
                excelData.push(['To Date:', dateToFilter]);
            }
        }
        
        // Create workbook and worksheet
        const wb = XLSX.utils.book_new();
        const ws = XLSX.utils.aoa_to_sheet(excelData);
        
        // Set column widths
        const colWidths = [
            { wch: 12 }, // Date
            { wch: 12 }, // Weekday
            { wch: 25 }, // Task
            { wch: 15 }, // Category
            { wch: 35 }, // Description
            { wch: 12 }, // Duration (Minutes)
            { wch: 12 }, // Duration (Hours)
            { wch: 12 }, // Start Time
            { wch: 12 }  // End Time
        ];
        ws['!cols'] = colWidths;
        
        // Add worksheet to workbook
        XLSX.utils.book_append_sheet(wb, ws, 'Time Entries');
        
        // Generate filename with current date and filters
        let filename = 'timesheet_entries_' + new Date().toISOString().split('T')[0];
        if (categoryFilter) {
            filename += '_' + categoryFilter.replace(/\s+/g, '_');
        }
        if (dateFromFilter && dateToFilter) {
            filename += '_' + dateFromFilter + '_to_' + dateToFilter;
        } else if (dateFromFilter) {
            filename += '_from_' + dateFromFilter;
        } else if (dateToFilter) {
            filename += '_until_' + dateToFilter;
        }
        filename += '.xlsx';
        
        // Save file
        XLSX.writeFile(wb, filename);
        
        Utils.showSuccess(`Excel file exported: ${filename}`);
        
    } catch (error) {
        console.error('Error exporting to Excel:', error);
        Utils.showError('Failed to export Excel file. Please try again.');
    }
}