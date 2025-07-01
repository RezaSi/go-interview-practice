// Common JavaScript utilities for Go Interview Practice web UI

// Helper for formatting timestamps
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Helper for truncating text
function truncateText(text, maxLength = 150) {
    if (text.length <= maxLength) return text;
    return text.substr(0, maxLength) + '...';
}

// Helper for formatting execution time
function formatExecutionTime(ms) {
    if (ms < 1000) {
        return ms + 'ms';
    } else {
        return (ms / 1000).toFixed(2) + 's';
    }
}

// Highlighting functionality for learning materials
class LearningHighlighter {
    constructor(containerId, challengeId) {
        this.container = document.getElementById(containerId);
        this.challengeId = challengeId;
        this.storageKey = `highlights_challenge_${challengeId}`;
        this.isHighlighting = false;
        this.highlightColor = '#ffeb3b'; // Default yellow
        this.highlights = this.loadHighlights();
        this.init();
    }

    init() {
        if (!this.container) return;
        
        // Add highlighting controls
        this.addControls();
        
        // Apply existing highlights
        this.applyHighlights();
        
        // Add event listeners
        this.addEventListeners();
    }

    addControls() {
        const controlsDiv = document.createElement('div');
        controlsDiv.className = 'highlighting-controls mb-3 p-2 bg-light rounded border';
        controlsDiv.innerHTML = `
            <div class="d-flex align-items-center justify-content-between flex-wrap gap-2">
                <div class="d-flex align-items-center gap-2">
                    <span class="fw-semibold text-muted small">Highlights:</span>
                    <button class="btn btn-sm btn-outline-primary" id="toggle-highlight-btn">
                        <i class="bi bi-highlighter"></i> Toggle
                    </button>
                    <div class="d-flex align-items-center gap-1">
                        <span class="text-muted small">Color:</span>
                        <div class="color-options">
                            <button class="color-btn active" data-color="#ffeb3b" style="background-color: #ffeb3b;" title="Yellow"></button>
                            <button class="color-btn" data-color="#ffcdd2" style="background-color: #ffcdd2;" title="Pink"></button>
                            <button class="color-btn" data-color="#c8e6c9" style="background-color: #c8e6c9;" title="Green"></button>
                            <button class="color-btn" data-color="#bbdefb" style="background-color: #bbdefb;" title="Blue"></button>
                        </div>
                    </div>
                </div>
                <div class="d-flex align-items-center gap-1">
                    <button class="btn btn-sm btn-outline-secondary" id="clear-highlights-btn" title="Clear all highlights">
                        <i class="bi bi-trash"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-info" id="export-highlights-btn" title="Export highlights">
                        <i class="bi bi-download"></i>
                    </button>
                </div>
            </div>
            <div class="mt-2 d-flex justify-content-between align-items-center">
                <small class="text-muted">
                    <span id="highlight-status">Mode: OFF</span>
                </small>
                <small class="text-muted">
                    <span id="highlight-count">${this.highlights.length} saved</span>
                </small>
            </div>
        `;
        
        this.container.insertBefore(controlsDiv, this.container.firstChild);
    }

    addEventListeners() {
        // Toggle highlight mode
        const toggleBtn = document.getElementById('toggle-highlight-btn');
        if (toggleBtn) {
            toggleBtn.addEventListener('click', () => this.toggleHighlightMode());
        }

        // Color buttons
        const colorBtns = document.querySelectorAll('.color-btn');
        colorBtns.forEach(btn => {
            btn.addEventListener('click', (e) => {
                // Remove active class from all buttons
                colorBtns.forEach(b => b.classList.remove('active'));
                // Add active class to clicked button
                e.target.classList.add('active');
                // Update highlight color
                this.highlightColor = e.target.dataset.color;
            });
        });

        // Clear all highlights
        const clearBtn = document.getElementById('clear-highlights-btn');
        if (clearBtn) {
            clearBtn.addEventListener('click', () => this.clearAllHighlights());
        }

        // Export highlights
        const exportBtn = document.getElementById('export-highlights-btn');
        if (exportBtn) {
            exportBtn.addEventListener('click', () => this.exportHighlights());
        }

        // Mouse events for highlighting
        this.container.addEventListener('mouseup', (e) => this.handleMouseUp(e));
    }

    toggleHighlightMode() {
        this.isHighlighting = !this.isHighlighting;
        const toggleBtn = document.getElementById('toggle-highlight-btn');
        const status = document.getElementById('highlight-status');
        
        if (this.isHighlighting) {
            toggleBtn.classList.remove('btn-outline-primary');
            toggleBtn.classList.add('btn-primary');
            toggleBtn.innerHTML = '<i class="bi bi-highlighter"></i> ON';
            status.textContent = 'Mode: ON - Click & drag to highlight';
            this.container.style.cursor = 'crosshair';
        } else {
            toggleBtn.classList.remove('btn-primary');
            toggleBtn.classList.add('btn-outline-primary');
            toggleBtn.innerHTML = '<i class="bi bi-highlighter"></i> Toggle';
            status.textContent = 'Mode: OFF';
            this.container.style.cursor = 'default';
        }
    }

    handleMouseUp(e) {
        if (!this.isHighlighting) return;
        
        const selection = window.getSelection();
        if (selection.toString().trim()) {
            this.highlightSelection(selection);
        }
    }

    highlightSelection(selection) {
        const range = selection.getRangeAt(0);
        const text = selection.toString().trim();
        
        if (!text) return;

        // Create highlight element
        const highlightSpan = document.createElement('span');
        highlightSpan.className = 'highlighted-text';
        highlightSpan.style.backgroundColor = this.highlightColor;
        highlightSpan.style.cursor = 'pointer';
        highlightSpan.title = 'Click to remove highlight';
        
        // Add delete button
        const deleteBtn = document.createElement('span');
        deleteBtn.className = 'highlight-delete-btn';
        deleteBtn.innerHTML = '×';
        deleteBtn.title = 'Remove highlight';
        deleteBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            this.removeHighlight(highlightSpan);
        });
        
        // Extract the text content and create a new range
        const fragment = range.extractContents();
        highlightSpan.appendChild(fragment);
        highlightSpan.appendChild(deleteBtn);
        range.insertNode(highlightSpan);
        
        // Save highlight data
        const highlightData = {
            id: Date.now() + Math.random(),
            text: text,
            color: this.highlightColor,
            timestamp: new Date().toISOString(),
            context: this.getTextContext(text, 50)
        };
        
        this.highlights.push(highlightData);
        this.saveHighlights();
        this.updateHighlightCount();
        
        // Clear selection
        selection.removeAllRanges();
    }

    removeHighlight(element) {
        // Find and remove from highlights array
        const text = element.textContent.replace('×', '').trim(); // Remove the × character
        this.highlights = this.highlights.filter(h => h.text !== text);
        this.saveHighlights();
        this.updateHighlightCount();
        
        // Remove from DOM
        const parent = element.parentNode;
        parent.replaceChild(document.createTextNode(text), element);
        parent.normalize(); // Merge adjacent text nodes
    }

    getTextContext(text, contextLength) {
        const container = this.container;
        const fullText = container.textContent;
        const textIndex = fullText.indexOf(text);
        
        if (textIndex === -1) return '';
        
        const start = Math.max(0, textIndex - contextLength);
        const end = Math.min(fullText.length, textIndex + text.length + contextLength);
        
        return fullText.substring(start, end).replace(/\s+/g, ' ').trim();
    }

    applyHighlights() {
        this.highlights.forEach(highlight => {
            this.applyHighlight(highlight);
        });
    }

    applyHighlight(highlightData) {
        const walker = document.createTreeWalker(
            this.container,
            NodeFilter.SHOW_TEXT,
            null,
            false
        );

        const textNodes = [];
        let node;
        while (node = walker.nextNode()) {
            textNodes.push(node);
        }

        textNodes.forEach(textNode => {
            const text = textNode.textContent;
            const index = text.indexOf(highlightData.text);
            
            if (index !== -1) {
                const before = text.substring(0, index);
                const highlighted = text.substring(index, index + highlightData.text.length);
                const after = text.substring(index + highlightData.text.length);
                
                const fragment = document.createDocumentFragment();
                
                if (before) {
                    fragment.appendChild(document.createTextNode(before));
                }
                
                const highlightSpan = document.createElement('span');
                highlightSpan.className = 'highlighted-text';
                highlightSpan.style.backgroundColor = highlightData.color;
                highlightSpan.style.cursor = 'pointer';
                highlightSpan.title = 'Click to remove highlight';
                highlightSpan.textContent = highlighted;
                
                // Add delete button
                const deleteBtn = document.createElement('span');
                deleteBtn.className = 'highlight-delete-btn';
                deleteBtn.innerHTML = '×';
                deleteBtn.title = 'Remove highlight';
                deleteBtn.addEventListener('click', (e) => {
                    e.stopPropagation();
                    this.removeHighlight(highlightSpan);
                });
                
                highlightSpan.appendChild(deleteBtn);
                fragment.appendChild(highlightSpan);
                
                if (after) {
                    fragment.appendChild(document.createTextNode(after));
                }
                
                textNode.parentNode.replaceChild(fragment, textNode);
            }
        });
    }

    clearAllHighlights() {
        if (confirm('Are you sure you want to clear all highlights?')) {
            this.highlights = [];
            this.saveHighlights();
            this.updateHighlightCount();
            
            // Remove all highlight spans from DOM
            const highlights = this.container.querySelectorAll('.highlighted-text');
            highlights.forEach(highlight => {
                const parent = highlight.parentNode;
                parent.replaceChild(document.createTextNode(highlight.textContent), highlight);
            });
            
            // Normalize text nodes
            this.container.normalize();
        }
    }

    exportHighlights() {
        if (this.highlights.length === 0) {
            alert('No highlights to export');
            return;
        }
        
        const exportData = {
            challengeId: this.challengeId,
            exportDate: new Date().toISOString(),
            highlights: this.highlights
        };
        
        const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `highlights_challenge_${this.challengeId}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    loadHighlights() {
        try {
            const saved = localStorage.getItem(this.storageKey);
            return saved ? JSON.parse(saved) : [];
        } catch (e) {
            console.error('Error loading highlights:', e);
            return [];
        }
    }

    saveHighlights() {
        try {
            localStorage.setItem(this.storageKey, JSON.stringify(this.highlights));
        } catch (e) {
            console.error('Error saving highlights:', e);
        }
    }

    updateHighlightCount() {
        const countElement = document.getElementById('highlight-count');
        if (countElement) {
            countElement.textContent = `${this.highlights.length} highlights saved`;
        }
    }
}

// Initialize syntax highlighting for code blocks
function initSyntaxHighlighting() {
    document.querySelectorAll('pre code').forEach((el) => {
        // Fix for Go language blocks
        if (el.className === 'language-go') {
            el.className = 'language-golang'; // Convert 'go' to 'golang' for better highlighting
        }
        hljs.highlightElement(el);
    });
}

// Initialize Markdown parsing
function renderMarkdown(markdownText, targetElement) {
    if (!markdownText || !targetElement) return;
    
    // Set options for the markdown parser
    marked.setOptions({
        gfm: true,      // GitHub flavored markdown
        breaks: true,   // Line breaks are rendered
        highlight: function(code, lang) {
            // Normalize language name - convert 'go' to 'golang' if needed
            let language = lang;
            if (lang === 'go') {
                language = 'golang';
            } else {
                language = hljs.getLanguage(lang) ? lang : 'plaintext';
            }
            return hljs.highlight(code, { language }).value;
        },
        langPrefix: 'language-'  // Add proper language prefix for CSS
    });
    
    // Parse and render markdown
    targetElement.innerHTML = marked.parse(markdownText);
    
    // Additionally, ensure Go code blocks are properly highlighted
    targetElement.querySelectorAll('pre code.language-go').forEach((el) => {
        el.className = 'language-golang';
        hljs.highlightElement(el);
    });
}

// Initialize Markdown parsing with cleanup (for challenge descriptions)
function renderMarkdownAndCleanup(markdownText, targetElement) {
    if (!markdownText || !targetElement) return;
    
    // Set options for the markdown parser
    marked.setOptions({
        gfm: true,      // GitHub flavored markdown
        breaks: true,   // Line breaks are rendered
        highlight: function(code, lang) {
            // Normalize language name - convert 'go' to 'golang' if needed
            let language = lang;
            if (lang === 'go') {
                language = 'golang';
            } else {
                language = hljs.getLanguage(lang) ? lang : 'plaintext';
            }
            return hljs.highlight(code, { language }).value;
        },
        langPrefix: 'language-'  // Add proper language prefix for CSS
    });
    
    // Parse and render markdown
    targetElement.innerHTML = marked.parse(markdownText);
    
    // Clean up any script tags for security
    const scripts = targetElement.querySelectorAll('script');
    scripts.forEach(script => script.remove());
    
    // Additionally, ensure Go code blocks are properly highlighted
    targetElement.querySelectorAll('pre code.language-go').forEach((el) => {
        el.className = 'language-golang';
        hljs.highlightElement(el);
    });
}

// Escape HTML for safe output
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Helper for creating a code editor
function createEditor(elementId, code, isReadOnly = false) {
    const editor = ace.edit(elementId);
    editor.setTheme("ace/theme/chrome");
    editor.session.setMode("ace/mode/golang");
    editor.setValue(code || '');
    editor.setReadOnly(isReadOnly);
    editor.clearSelection();
    editor.setOptions({
        fontSize: "14px",
        showPrintMargin: false,
        showGutter: true,
        highlightActiveLine: true,
        wrap: true
    });
    return editor;
}

// Save editor content in localStorage to prevent loss on page refresh
function saveEditorContent(key, content) {
    localStorage.setItem(`editor_${key}`, content);
}

// Load editor content from localStorage if available
function loadEditorContent(key) {
    return localStorage.getItem(`editor_${key}`);
}

// Format the test output for display
function formatTestOutput(output) {
    if (!output) return '';
    
    // Colorize PASS/FAIL in the output
    return output
        .replace(/PASS/g, '<span class="text-success">PASS</span>')
        .replace(/FAIL/g, '<span class="text-danger">FAIL</span>')
        .replace(/--- FAIL/g, '<span class="text-danger">--- FAIL</span>')
        .replace(/--- PASS/g, '<span class="text-success">--- PASS</span>');
}

// Handle form submissions with AJAX
function handleFormSubmit(formElement, successCallback, errorCallback) {
    formElement.addEventListener('submit', function(e) {
        e.preventDefault();
        const formData = new FormData(formElement);
        const jsonData = {};
        
        for (const [key, value] of formData.entries()) {
            jsonData[key] = value;
        }
        
        fetch(formElement.action, {
            method: formElement.method || 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(jsonData)
        })
        .then(response => response.json())
        .then(data => {
            if (typeof successCallback === 'function') {
                successCallback(data);
            }
        })
        .catch(error => {
            if (typeof errorCallback === 'function') {
                errorCallback(error);
            } else {
                console.error('Error:', error);
            }
        });
    });
}

// Custom template function for truncating text in templates
function truncateDescription(text, maxLength = 100) {
    if (!text) return '';
    const plainText = text.replace(/<[^>]*>/g, '');
    if (plainText.length <= maxLength) return plainText;
    return plainText.substr(0, maxLength) + '...';
}

// Initialize common page elements
document.addEventListener('DOMContentLoaded', function() {
    // Initialize tooltips
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
    
    // Initialize syntax highlighting
    initSyntaxHighlighting();
    
    // Apply renderMarkdown to all markdown content containers
    document.querySelectorAll('.markdown-content').forEach(function(el) {
        // If there's a data attribute with the markdown content, use that
        const content = el.getAttribute('data-markdown');
        if (content) {
            renderMarkdown(content, el);
        }
        // Otherwise just ensure proper syntax highlighting for existing content
        else {
            // Make sure Go code blocks have the right class for syntax highlighting
            el.querySelectorAll('pre code.language-go').forEach(function(codeEl) {
                codeEl.className = 'language-golang';
                hljs.highlightElement(codeEl);
            });
        }
    });
    
    // Handle username persistence
    const usernameInput = document.getElementById('username');
    if (usernameInput) {
        // Load saved username
        const savedUsername = localStorage.getItem('githubUsername');
        if (savedUsername) {
            usernameInput.value = savedUsername;
        }
        
        // Save username when changed
        usernameInput.addEventListener('change', function() {
            localStorage.setItem('githubUsername', this.value);
        });
    }
});

// Initialize learning materials with highlighting
function initLearningMaterials(containerId, challengeId) {
    const container = document.getElementById(containerId);
    if (container) {
        // Wait a bit for markdown to be rendered
        setTimeout(() => {
            new LearningHighlighter(containerId, challengeId);
        }, 100);
    }
}

// Initialize hints system
function initializeHints(hintsText) {
    if (!hintsText) return;
    
    try {
        const hints = JSON.parse(hintsText);
        const container = document.getElementById('hints-container');
        const showHintBtn = document.getElementById('show-hint-btn');
        const resetHintsBtn = document.getElementById('reset-hints-btn');
        const progressSpan = document.getElementById('hints-progress');
        const totalHintsSpan = document.getElementById('total-hints');
        
        if (!container || !showHintBtn) return;
        
        let currentHintIndex = 0;
        const storageKey = `hints_progress_${window.location.pathname.split('/').pop()}`;
        
        // Load saved progress
        const savedProgress = localStorage.getItem(storageKey);
        if (savedProgress) {
            currentHintIndex = parseInt(savedProgress);
        }
        
        // Update total hints count
        totalHintsSpan.textContent = hints.length;
        
        // Show hints up to current index
        for (let i = 0; i < currentHintIndex; i++) {
            if (hints[i]) {
                this.showHint(hints[i], i + 1);
            }
        }
        
        // Update progress display
        this.updateHintsProgress();
        
        // Show hint button click handler
        showHintBtn.addEventListener('click', () => {
            if (currentHintIndex < hints.length) {
                this.showHint(hints[currentHintIndex], currentHintIndex + 1);
                currentHintIndex++;
                localStorage.setItem(storageKey, currentHintIndex.toString());
                this.updateHintsProgress();
                
                if (currentHintIndex >= hints.length) {
                    showHintBtn.classList.add('d-none');
                    resetHintsBtn.classList.remove('d-none');
                }
            }
        });
        
        // Reset hints button click handler
        resetHintsBtn.addEventListener('click', () => {
            if (confirm('Are you sure you want to reset all hints?')) {
                container.innerHTML = '';
                currentHintIndex = 0;
                localStorage.removeItem(storageKey);
                this.updateHintsProgress();
                showHintBtn.classList.remove('d-none');
                resetHintsBtn.classList.add('d-none');
            }
        });
        
        // Helper function to show a hint
        this.showHint = function(hint, hintNumber) {
            const hintDiv = document.createElement('div');
            hintDiv.className = 'hint-item alert alert-warning';
            hintDiv.innerHTML = `
                <div class="d-flex justify-content-between align-items-start">
                    <div class="flex-grow-1">
                        <h6 class="alert-heading">
                            <i class="bi bi-lightbulb me-2"></i>Hint ${hintNumber}
                        </h6>
                        <div class="markdown-content">${marked.parse(hint)}</div>
                    </div>
                    <span class="badge bg-warning text-dark ms-2">${hintNumber}</span>
                </div>
            `;
            container.appendChild(hintDiv);
            
            // Apply syntax highlighting to code blocks in the hint
            setTimeout(() => {
                hintDiv.querySelectorAll('pre code').forEach((el) => {
                    hljs.highlightElement(el);
                });
            }, 100);
        };
        
        // Helper function to update progress display
        this.updateHintsProgress = function() {
            progressSpan.textContent = currentHintIndex;
        };
        
    } catch (error) {
        console.error('Error initializing hints:', error);
    }
} 