document.addEventListener('DOMContentLoaded', () => {

    const validPages = ['Home', 'Error', 'Categories'];
    let isAuthenticated = false;

    loadUserState();
    Authenticated(); 
    setupNavigationListeners();
    window.navigateToPage = navigateToPage;
    window.Authenticated = Authenticated;

    function Authenticated() {
        
        fetch("/auth/status", {
            method: "GET",
            credentials: "same-origin"  // Important to include cookies in the request
        })
        .then(response => response.json())
        .then(data => {
            const loginButtons = document.querySelectorAll(".login-buttons");
            const logoutButtons = document.querySelectorAll(".logout-buttons");
            const privilage = document.querySelectorAll(".privilage");
    
            const mainContent = document.querySelector('.main-content');
            const RightBar = document.querySelector('.RightSidebar')
    
            
                if (data.authenticated) {
                    isAuthenticated = true;
                    loginButtons.forEach(section => section.classList.add('deactive'));
                    logoutButtons.forEach(section => section.classList.remove('deactive'));
                    privilage.forEach(section => section.classList.remove('disabled-link'));
                    
                    mainContent.classList.remove('RightBar');
                    RightBar.classList.remove('deactive');

                    addValidPages(['Createpost','Profile','Created','Liked','Disliked','Messages']);
                    StartWebSocket();
                    
                    // if (data.privilege === 1){} else if (data.privilege === 2) {} else if (data.privilege === 3) {}
                } else {
                    isAuthenticated = false;
                    loginButtons.forEach(section => section.classList.remove('deactive'));
                    logoutButtons.forEach(section => section.classList.add('deactive'));
                    privilage.forEach(section => section.classList.add('disabled-link'));
                    
                    addValidPages(['Login', 'Register']);
                    mainContent.classList.add('RightBar');
                    RightBar.classList.add('deactive');
                }

                handleInitialPageLoad();
        })
        .catch(error => {
            console.error("Error checking auth status:", error)
            navigateToPage('Error');
        });

    }
    
    function setupNavigationListeners() {
        // Add event listeners for navigation buttons
        document.querySelectorAll('.sidebar [data-page], .main-content [data-page]').forEach(element => {
            element.addEventListener('click', debounce(() => {
                const pageId = element.dataset.page;
                navigateToPage(pageId);
            }, 200));
        });

        // Handle the browser's back/forward navigation
        window.addEventListener('popstate', () => {
            const currentPath = window.location.pathname.slice(1);
            const pageId = getValidPageId(currentPath);
            showPage(pageId);
        });
    }

    function handleInitialPageLoad() {
        const currentPath = window.location.pathname.slice(1);
        const pageIdToShow = getValidPageId(currentPath) || 'Error';
        showPage(pageIdToShow);
    }

    function navigateToPage(pageId) {
        if (validPages.includes(pageId)) {
            const path = pageId === 'Home' ? '/' : `/${pageId.toLowerCase()}`;
            history.pushState({}, '', path);
            showPage(pageId);
            localStorage.setItem('currentPage', pageId);
        } else {
            if (isAuthenticated) {
                console.error(`Invalid page ID: ${pageId}`);
                navigateToPage('Error');
            } else {
                navigateToPage('Login');
            }
        }
    }
    // Show the requested page
    function showPage(page) {
        const contents = document.querySelectorAll('.deactive');
        contents.forEach(content => {
            content.classList.remove('active');
        });
        
        const activeContent = document.getElementById(page);
        if (activeContent) {
            activeContent.classList.add('active');
        }
        if (page === 'Home') {
            loadPosts();
            console.log("Posts loaded, applying handlers...");
        } else if (page === 'Categories') {
            loadCategories();
            console.log("Categories loaded, applying handlers...");
        } else if (page === 'Profile') {
            loadProfileData();
            console.log("Categories loaded, applying handlers...");
        }
    }
    
    function addValidPages(pages) {
        pages.forEach(page => {
            if (!validPages.includes(page)) {
                validPages.push(page);
            }
        });
    }

    function getValidPageId(path) {
        const pageId = path === '' ? 'Home' : path.charAt(0).toUpperCase() + path.slice(1);
        return validPages.includes(pageId) ? pageId : null;
    }
});