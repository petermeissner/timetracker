/**
 * Centralized API access layer for the timesheet application
 * Provides consistent error handling, request formatting, and response processing
 */

const API = {
    /**
     * Base configuration and utilities
     */
    baseURL: '/api',
    
    /**
     * Generic request handler with consistent error handling
     */
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const defaultHeaders = {
            'Content-Type': 'application/json',
            ...options.headers
        };
        
        try {
            const response = await fetch(url, {
                ...options,
                headers: defaultHeaders
            });
            
            if (!response.ok) {
                let errorMessage;
                try {
                    errorMessage = await response.text();
                } catch {
                    errorMessage = `HTTP ${response.status}: ${response.statusText}`;
                }
                throw new Error(`Server error (${response.status}): ${errorMessage}`);
            }
            
            // Handle empty responses (like DELETE operations)
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return await response.json();
            }
            
            return null;
        } catch (error) {
            // Re-throw with additional context
            console.error(`API Request failed: ${options.method || 'GET'} ${url}`, error);
            throw error;
        }
    },
    
    /**
     * Categories API
     */
    categories: {
        // GET /api/categories
        async getAll() {
            return API.request('/categories');
        },
        
        // POST /api/categories
        async create(categoryData) {
            return API.request('/categories', {
                method: 'POST',
                body: JSON.stringify(categoryData)
            });
        },
        
        // PUT /api/categories/:id
        async update(id, categoryData) {
            return API.request(`/categories/${id}`, {
                method: 'PUT',
                body: JSON.stringify(categoryData)
            });
        },
        
        // DELETE /api/categories/:id
        async delete(id) {
            return API.request(`/categories/${id}`, {
                method: 'DELETE'
            });
        }
    },
    
    /**
     * Tasks API
     */
    tasks: {
        // GET /api/tasks
        async getAll() {
            return API.request('/tasks');
        },
        
        // POST /api/tasks
        async create(taskData) {
            return API.request('/tasks', {
                method: 'POST',
                body: JSON.stringify(taskData)
            });
        },
        
        // PUT /api/tasks/:id
        async update(id, taskData) {
            return API.request(`/tasks/${id}`, {
                method: 'PUT',
                body: JSON.stringify(taskData)
            });
        },
        
        // DELETE /api/tasks/:id
        async delete(id) {
            return API.request(`/tasks/${id}`, {
                method: 'DELETE'
            });
        }
    },
    
    /**
     * Time Entries API
     */
    entries: {
        // GET /api/entries
        async getAll() {
            return API.request('/entries');
        },
        
        // POST /api/entries
        async create(entryData) {
            return API.request('/entries', {
                method: 'POST',
                body: JSON.stringify(entryData)
            });
        },
        
        // PUT /api/entries/:id
        async update(id, entryData) {
            return API.request(`/entries/${id}`, {
                method: 'PUT',
                body: JSON.stringify(entryData)
            });
        },
        
        // DELETE /api/entries/:id
        async delete(id) {
            return API.request(`/entries/${id}`, {
                method: 'DELETE'
            });
        }
    }
};

// Export for use in other modules
window.API = API;