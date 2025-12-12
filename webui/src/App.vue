<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>

	<header class="navbar custom-navbar sticky-top">
     <a class="navbar-brand custom-brand" href="">WASAText</a>
     <button class="navbar-toggler custom-toggler" type="button">
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
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								My Conversations
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/me" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								My Profile
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/search" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
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
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
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

/* Navbar toggle button */
.custom-toggler {
  background-color: #fff;
  border: none;
  width: 50px;
  height: 40px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.custom-toggler:hover {
  background-color: #e0e0e0;
}
</style>
