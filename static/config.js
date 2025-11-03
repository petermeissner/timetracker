// Global variables
let categories = [];
let tasks = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadCategories();
    loadTasks();
    setupEventListeners();
});

function setupEventListeners() {
    // Category events
    document.getElementById('addCategoryBtn').addEventListener('click', showAddCategoryModal);
    document.getElementById('categoryForm').addEventListener('submit', handleCategorySubmit);
    document.getElementById('categoryColor').addEventListener('input', updateColorPreview);
    
    // Task events
    document.getElementById('addTaskBtn').addEventListener('click', showAddTaskModal);
    document.getElementById('taskForm').addEventListener('submit', handleTaskSubmit);
    
    // Modal close events
    window.addEventListener('click', function(event) {
        const categoryModal = document.getElementById('categoryModal');
        const taskModal = document.getElementById('taskModal');
        if (event.target === categoryModal) {
            closeCategoryModal();
        }
        if (event.target === taskModal) {
            closeTaskModal();
        }
    });
}

// Category functions
async function loadCategories() {
    try {
        const response = await fetch('/api/categories');
        if (!response.ok) {
            throw new Error('Failed to load categories');
        }
        
        categories = await response.json() || [];
        renderCategories();
        populateTaskCategorySelect();
    } catch (error) {
        console.error('Error loading categories:', error);
        showError('Failed to load categories');
    }
}

function renderCategories() {
    const categoriesList = document.getElementById('categoriesList');
    
    if (categories.length === 0) {
        categoriesList.innerHTML = `
            <div class="empty-state">
                <h3>No categories configured</h3>
                <p>Add your first category to get started!</p>
            </div>
        `;
        return;
    }
    
    let html = '';
    categories.forEach(category => {
        html += `
            <div class="config-item" data-id="${category.id}">
                <div class="config-item-header">
                    <div>
                        <div class="config-item-name">
                            <span class="category-color-indicator" style="background-color: ${category.color}"></span>
                            ${escapeHtml(category.name)}
                        </div>
                        <div class="config-item-details">
                            Color: ${category.color}
                        </div>
                    </div>
                    <div class="config-item-actions">
                        <button class="btn btn-secondary btn-small" onclick="editCategory(${category.id})">Edit</button>
                        <button class="btn btn-danger btn-small" onclick="deleteCategory(${category.id})">Delete</button>
                    </div>
                </div>
            </div>
        `;
    });
    
    categoriesList.innerHTML = html;
}

function showAddCategoryModal() {
    document.getElementById('categoryModalTitle').textContent = 'Add Category';
    document.getElementById('categoryForm').reset();
    document.getElementById('categoryId').value = '';
    document.getElementById('categoryColor').value = '#718096';
    updateColorPreview();
    document.getElementById('categoryModal').style.display = 'block';
}

function editCategory(id) {
    const category = categories.find(c => c.id === id);
    if (!category) return;
    
    document.getElementById('categoryModalTitle').textContent = 'Edit Category';
    document.getElementById('categoryId').value = category.id;
    document.getElementById('categoryName').value = category.name;
    document.getElementById('categoryColor').value = category.color;
    updateColorPreview();
    document.getElementById('categoryModal').style.display = 'block';
}

async function handleCategorySubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const categoryId = document.getElementById('categoryId').value;
    const data = {
        name: formData.get('name').trim(),
        color: formData.get('color')
    };
    
    if (!data.name) {
        showError('Category name is required');
        return;
    }
    
    try {
        const url = categoryId ? `/api/categories/${categoryId}` : '/api/categories';
        const method = categoryId ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error('Failed to save category');
        }
        
        await loadCategories();
        closeCategoryModal();
        showSuccess(`Category ${categoryId ? 'updated' : 'created'} successfully!`);
    } catch (error) {
        console.error('Error saving category:', error);
        showError('Failed to save category');
    }
}

async function deleteCategory(id) {
    if (!confirm('Are you sure you want to delete this category? This action cannot be undone.')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/categories/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete category');
        }
        
        await loadCategories();
        showSuccess('Category deleted successfully!');
    } catch (error) {
        console.error('Error deleting category:', error);
        showError('Failed to delete category');
    }
}

function closeCategoryModal() {
    document.getElementById('categoryModal').style.display = 'none';
    document.getElementById('categoryForm').reset();
}

function updateColorPreview() {
    const color = document.getElementById('categoryColor').value;
    const preview = document.querySelector('.color-preview');
    if (preview) {
        preview.style.backgroundColor = color;
    }
}

// Task functions
async function loadTasks() {
    try {
        const response = await fetch('/api/tasks');
        if (!response.ok) {
            throw new Error('Failed to load tasks');
        }
        
        tasks = await response.json() || [];
        renderTasks();
    } catch (error) {
        console.error('Error loading tasks:', error);
        showError('Failed to load tasks');
    }
}

function renderTasks() {
    const tasksList = document.getElementById('tasksList');
    
    if (tasks.length === 0) {
        tasksList.innerHTML = `
            <div class="empty-state">
                <h3>No predefined tasks</h3>
                <p>Add predefined tasks to make time entry faster!</p>
            </div>
        `;
        return;
    }
    
    let html = '';
    tasks.forEach(task => {
        const category = categories.find(c => c.id === task.category_id);
        const categoryInfo = category ? {
            name: category.name,
            color: category.color
        } : null;
        
        html += `
            <div class="config-item" data-id="${task.id}">
                <div class="config-item-header">
                    <div>
                        <div class="config-item-name">${escapeHtml(task.name)}</div>
                        <div class="config-item-details">
                            ${task.description ? escapeHtml(task.description) : 'No description'}
                        </div>
                        ${categoryInfo ? `
                            <div class="task-category-info">
                                <span class="task-category-badge" style="background-color: ${categoryInfo.color}">
                                    ${escapeHtml(categoryInfo.name)}
                                </span>
                            </div>
                        ` : ''}
                    </div>
                    <div class="config-item-actions">
                        <button class="btn btn-secondary btn-small" onclick="editTask(${task.id})">Edit</button>
                        <button class="btn btn-danger btn-small" onclick="deleteTask(${task.id})">Delete</button>
                    </div>
                </div>
            </div>
        `;
    });
    
    tasksList.innerHTML = html;
}

function populateTaskCategorySelect() {
    const select = document.getElementById('taskCategory');
    let html = '<option value="">No Category</option>';
    
    categories.forEach(category => {
        html += `<option value="${category.id}">${escapeHtml(category.name)}</option>`;
    });
    
    select.innerHTML = html;
}

function showAddTaskModal() {
    document.getElementById('taskModalTitle').textContent = 'Add Task';
    document.getElementById('taskForm').reset();
    document.getElementById('taskId').value = '';
    document.getElementById('taskModal').style.display = 'block';
}

function editTask(id) {
    const task = tasks.find(t => t.id === id);
    if (!task) return;
    
    document.getElementById('taskModalTitle').textContent = 'Edit Task';
    document.getElementById('taskId').value = task.id;
    document.getElementById('taskName').value = task.name;
    document.getElementById('taskCategory').value = task.category_id || '';
    document.getElementById('taskDescription').value = task.description || '';
    document.getElementById('taskModal').style.display = 'block';
}

async function handleTaskSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const taskId = document.getElementById('taskId').value;
    const data = {
        name: formData.get('name').trim(),
        category_id: parseInt(formData.get('category_id')) || 0,
        description: formData.get('description').trim()
    };
    
    if (!data.name) {
        showError('Task name is required');
        return;
    }
    
    try {
        const url = taskId ? `/api/tasks/${taskId}` : '/api/tasks';
        const method = taskId ? 'PUT' : 'POST';
        
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            throw new Error('Failed to save task');
        }
        
        await loadTasks();
        closeTaskModal();
        showSuccess(`Task ${taskId ? 'updated' : 'created'} successfully!`);
    } catch (error) {
        console.error('Error saving task:', error);
        showError('Failed to save task');
    }
}

async function deleteTask(id) {
    if (!confirm('Are you sure you want to delete this task?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/tasks/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete task');
        }
        
        await loadTasks();
        showSuccess('Task deleted successfully!');
    } catch (error) {
        console.error('Error deleting task:', error);
        showError('Failed to delete task');
    }
}

function closeTaskModal() {
    document.getElementById('taskModal').style.display = 'none';
    document.getElementById('taskForm').reset();
}

// Utility functions
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