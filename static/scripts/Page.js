// Toggle sidebar collapse
let isCollapsed = false;
const sidebar = document.getElementById('sidebar');
const mainContent = document.querySelector('.main-content');
const usersidebar = document.querySelector('.users-sidebar');

document.addEventListener('DOMContentLoaded', () => {

    const storedState = localStorage.getItem('isCollapsed');
    if (storedState !== null) {
        isCollapsed = JSON.parse(storedState);
        sidebar.classList.toggle('collapsed', isCollapsed);
        mainContent.classList.toggle('collapsed', isCollapsed);
        usersidebar.classList.toggle('collapsed', isCollapsed);
    }

    const menuItems = document.querySelectorAll('.menu-item');

    menuItems.forEach(item => {
        item.addEventListener('click', (event) => {
            // Remove active class from all items
            menuItems.forEach(i => i.classList.remove('color'));
            // Add active class to clicked item
            item.classList.add('color');

            if (item.dataset.hasSubmenu) {
                const subMenu = item.nextElementSibling;
                const allSubMenus = document.querySelectorAll('.sub-menu');
                const arrowIcon = item.querySelector('i.material-icons:last-child');
                const allArrows = document.querySelectorAll('.menu-item i.material-icons:last-child');

                allSubMenus.forEach(menu => {
                    if (menu !== subMenu) {
                        menu.classList.remove('show');
                    }
                });

                allArrows.forEach(icon => {
                    if (icon !== arrowIcon) {
                        icon.textContent = 'arrow_drop_down';
                    }
                });
                subMenu.classList.toggle('show');
                arrowIcon.textContent = subMenu.classList.contains('show') ? 'arrow_drop_up' : 'arrow_drop_down';

            }
        });
        
    });

    document.querySelector('.list-menu').addEventListener('click', () => {
        isCollapsed = !isCollapsed;
        sidebar.classList.toggle('collapsed', isCollapsed);
        mainContent.classList.toggle('collapsed', isCollapsed);
        usersidebar.classList.toggle('collapsed', isCollapsed);

        localStorage.setItem('isCollapsed', JSON.stringify(isCollapsed));
    });

});