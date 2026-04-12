// Interactive Landing Page JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // Smooth scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });

    // Animate stats on scroll
    const observerOptions = {
        threshold: 0.5,
        rootMargin: '0px 0px -100px 0px'
    };

    const statsObserver = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateStats();
                statsObserver.unobserve(entry.target);
            }
        });
    }, observerOptions);

    const statsSection = document.querySelector('.stats-row');
    if (statsSection) {
        statsObserver.observe(statsSection);
    }

    function animateStats() {
        const statNumbers = document.querySelectorAll('.stat-number');
        statNumbers.forEach((stat, index) => {
            const finalValue = stat.textContent;
            const numericValue = parseInt(finalValue.replace(/\D/g, ''));
            const suffix = finalValue.replace(/\d/g, '');
            
            let current = 0;
            const increment = numericValue / 30;
            const timer = setInterval(() => {
                current += increment;
                if (current >= numericValue) {
                    stat.textContent = finalValue;
                    clearInterval(timer);
                } else {
                    stat.textContent = Math.floor(current) + suffix;
                }
            }, 50 + (index * 20));
        });
    }

    // Parallax effect for hero section
    window.addEventListener('scroll', () => {
        const scrolled = window.pageYOffset;
        const heroVisual = document.querySelector('.hero-visual');
        if (heroVisual) {
            heroVisual.style.transform = `translateY(${scrolled * 0.1}px)`;
        }
    });

    // Add loading animation to cards
    const cards = document.querySelectorAll('.feature-card, .category-card, .step-card');
    const cardObserver = new IntersectionObserver((entries) => {
        entries.forEach((entry, index) => {
            if (entry.isIntersecting) {
                setTimeout(() => {
                    entry.target.classList.add('loading');
                }, index * 100);
                cardObserver.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1 });

    cards.forEach(card => {
        cardObserver.observe(card);
    });

    // Code window typing effect
    const codeContent = document.querySelector('.go-code');
    if (codeContent) {
        const originalHTML = codeContent.innerHTML;
        codeContent.innerHTML = '';
        
        let i = 0;
        const typingSpeed = 15;
        
        // Strip HTML tags to get plain text for typing
        const tempDiv = document.createElement('div');
        tempDiv.innerHTML = originalHTML;
        const plainText = tempDiv.textContent || tempDiv.innerText || '';
        
        function typeCode() {
            if (i < plainText.length) {
                // Find the corresponding position in the HTML
                let currentHTML = '';
                let textIndex = 0;
                let htmlIndex = 0;
                
                while (textIndex <= i && htmlIndex < originalHTML.length) {
                    if (originalHTML[htmlIndex] === '<') {
                        // Skip HTML tag
                        while (htmlIndex < originalHTML.length && originalHTML[htmlIndex] !== '>') {
                            currentHTML += originalHTML[htmlIndex];
                            htmlIndex++;
                        }
                        if (htmlIndex < originalHTML.length) {
                            currentHTML += originalHTML[htmlIndex]; // Add the '>'
                            htmlIndex++;
                        }
                    } else {
                        currentHTML += originalHTML[htmlIndex];
                        if (originalHTML[htmlIndex] !== '\n' && originalHTML[htmlIndex] !== ' ' || textIndex === i) {
                            textIndex++;
                        }
                        htmlIndex++;
                    }
                }
                
                codeContent.innerHTML = currentHTML;
                i++;
                setTimeout(typeCode, typingSpeed);
            }
        }
        
        // Start typing when code window is visible
        const codeObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    setTimeout(typeCode, 500);
                    codeObserver.unobserve(entry.target);
                }
            });
        }, { threshold: 0.5 });
        
        codeObserver.observe(codeContent);
    }

    // Navbar background on scroll
    const navbar = document.querySelector('.navbar');
    window.addEventListener('scroll', () => {
        if (window.scrollY > 50) {
            navbar.style.background = 'rgba(255, 255, 255, 0.98)';
            navbar.style.boxShadow = '0 2px 20px rgba(0, 0, 0, 0.1)';
        } else {
            navbar.style.background = 'rgba(255, 255, 255, 0.95)';
            navbar.style.boxShadow = 'none';
        }
    });

    // Button click effects
    document.querySelectorAll('.btn').forEach(button => {
        button.addEventListener('click', function(e) {
            // Create ripple effect
            const ripple = document.createElement('span');
            const rect = this.getBoundingClientRect();
            const size = Math.max(rect.width, rect.height);
            const x = e.clientX - rect.left - size / 2;
            const y = e.clientY - rect.top - size / 2;
            
            ripple.style.width = ripple.style.height = size + 'px';
            ripple.style.left = x + 'px';
            ripple.style.top = y + 'px';
            ripple.classList.add('ripple');
            
            this.appendChild(ripple);
            
            setTimeout(() => {
                ripple.remove();
            }, 600);
        });
    });

    // Add ripple CSS
    const style = document.createElement('style');
    style.textContent = `
        .btn {
            position: relative;
            overflow: hidden;
        }
        .ripple {
            position: absolute;
            border-radius: 50%;
            background: rgba(255, 255, 255, 0.4);
            transform: scale(0);
            animation: ripple-animation 0.6s linear;
            pointer-events: none;
        }
        @keyframes ripple-animation {
            to {
                transform: scale(4);
                opacity: 0;
            }
        }
    `;
    document.head.appendChild(style);

    // Preload critical images
    const preloadImages = [
        // Add any image URLs you want to preload
    ];
    
    preloadImages.forEach(src => {
        const img = new Image();
        img.src = src;
    });

    // Performance optimization: Lazy load non-critical content
    if ('IntersectionObserver' in window) {
        const lazyElements = document.querySelectorAll('[data-lazy]');
        const lazyObserver = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const element = entry.target;
                    element.src = element.dataset.lazy;
                    element.removeAttribute('data-lazy');
                    lazyObserver.unobserve(element);
                }
            });
        });
        
        lazyElements.forEach(element => {
            lazyObserver.observe(element);
        });
    }

    // Fetch real-time GitHub stars
    async function fetchGitHubStars(githubUrl) {
        try {
            // Extract owner/repo from GitHub URL
            const parts = githubUrl.split('/');
            if (parts.length < 2) return null;
            
            const repo = `${parts[parts.length - 2]}/${parts[parts.length - 1]}`;
            const apiUrl = `https://api.github.com/repos/${repo}`;
            
            const response = await fetch(apiUrl);
            if (!response.ok) return null;
            
            const data = await response.json();
            return data.stargazers_count;
        } catch (error) {
            console.log('Error fetching GitHub stars:', error);
            return null;
        }
    }

    // Update all star counts on page
    async function updateStarCounts() {
        const starElements = document.querySelectorAll('[data-github-url]');
        
        for (const element of starElements) {
            const githubUrl = element.getAttribute('data-github-url');
            const stars = await fetchGitHubStars(githubUrl);
            
            if (stars !== null) {
                // Format stars (e.g., 35000 -> 35K)
                const formattedStars = stars >= 1000 ? 
                    Math.round(stars / 1000) + 'K' : 
                    stars.toString();
                
                element.textContent = formattedStars;
            }
        }
    }

    // Update star counts
    updateStarCounts();

    console.log('🚀 Go Interview Practice - Landing page loaded successfully!');
}); 