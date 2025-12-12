<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>

	<header class="navbar custom-navbar sticky-top">
     <a class="navbar-brand custom-brand" href="">WASAText</a>
     <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
       <span class="navbar-toggler-icon"></span>
     </button>
    </header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/home" class="nav-link">
							 My Conversations
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/me" class="nav-link">
							 My Profile
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/search" class="nav-link">
							 Search People
							</RouterLink>
						</li>
					</ul>

					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Groups</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/groups'" class="nav-link">
							 My Groups
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
/* Custom Navbar */
.custom-navbar {
  background-color: #4a90e2; /* blue background */
  height: 70px;
  padding: 0 20px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Navbar brand (WASAText) */
.custom-brand {
  font-size: 28px;       /* bigger font */
  font-weight: bold;     /* bold */
  color: #fff !important; /* white text */
  font-family: 'Segoe UI', sans-serif; /* custom font */
  letter-spacing: 2px;   /* spacing between letters */
}

/* Sidebar background */
.bg-custom-sidebar {
  background-color: #343a40; /* dark gray */
  color: #fff;
}

/* Sidebar links */
.sidebar-link {
  color: #fff;
  font-weight: 500;
  font-size: 16px;
  text-decoration: none;
}

.sidebar-link:hover {
  color: #2f36faff; /* gold hover */
}

/* Sidebar link default style */
.sidebar .nav-link {
  color: #fff; /* default text color */
  transition: color 0.3s; /* smooth color change */
}

/* Sidebar link on hover */
.sidebar .nav-link:hover {
  color: #2f36faff; /* change to any color you like on hover */
  text-decoration: none; /* remove underline if any */
}

/* Sidebar headings */
.sidebar-heading {
  color: #ccc;
  font-weight: bold;
  font-size: 14px;
}

/* Remove bullet points */
.nav.flex-column {
  list-style: none;
  padding-left: 0;
}

/* Adjust main content */
main {
  background-color: #f8f9fa; /* optional */
  min-height: 100vh;
  padding-top: 1rem;
}


</style>
