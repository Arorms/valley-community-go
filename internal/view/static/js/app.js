(function () {
  var root = document.documentElement;

  var themeToggle = document.getElementById('theme-toggle');
  if (themeToggle) {
    themeToggle.addEventListener('click', function () {
      root.classList.toggle('dark');
      try {
        localStorage.setItem('theme', root.classList.contains('dark') ? 'dark' : 'light');
      } catch (e) {}
    });
  }

  var sidebar = document.getElementById('sidebar');
  var overlay = document.getElementById('sidebar-overlay');
  var menuToggle = document.getElementById('menu-toggle');

  function closeSidebar() {
    if (sidebar) sidebar.classList.add('-translate-x-full');
    if (overlay) overlay.classList.add('hidden');
  }

  if (menuToggle) {
    menuToggle.addEventListener('click', function () {
      if (!sidebar) return;
      sidebar.classList.toggle('-translate-x-full');
      if (overlay) overlay.classList.toggle('hidden');
    });
  }
  if (overlay) {
    overlay.addEventListener('click', closeSidebar);
  }

  document.addEventListener('keydown', function (event) {
    if (event.key === 'Escape') closeSidebar();
  });
})();
