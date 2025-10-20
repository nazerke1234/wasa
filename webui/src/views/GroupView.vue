<script>
export default {
  data() {
    return {
      userDisplayName: "",
      errorNotification: null,
      isLoading: false,
      userGroups: []
    };
  },
  methods: {
    async fetchUserGroups() {
      this.errorNotification = null;
      this.isLoading = true;
      try {
        const authToken = localStorage.getItem("token");
        if (!authToken) {
          this.$router.push({ path: "/" });
          return;
        }
        const apiResponse = await this.$axios.get("/groups", {
          headers: {
            Authorization: `Bearer ${authToken}`
          }
        });
        this.userGroups = apiResponse.data || [];
      } catch (apiError) {
        console.error("Error fetching groups:", apiError);
        this.errorNotification = "Unable to load groups. Please refresh the page.";
      } finally {
        this.isLoading = false;
      }
    },
    navigateToGroup(groupIdentifier, groupTitle) {
        localStorage.setItem("groupName", groupTitle);
        this.$router.push({
            path: `/groups/${groupIdentifier}`
        });
    },
    reloadGroups() {
        this.fetchUserGroups();
    },
    signOut() {
        this.$router.push({ path: "/" });
    },
    createNewGroup() {
        this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    this.userDisplayName = localStorage.getItem("name") || "User";
    this.fetchUserGroups();
  }
};
</script>

<template>
  <div>
    <div class="header-container">
      <h1 class="page-heading">{{ userDisplayName }}, your group list</h1>
      <div class="action-buttons">
        <div class="button-group">
          <button type="button" class="secondary-button" @click="reloadGroups">Reload</button>
          <button type="button" class="secondary-button" @click="signOut">Sign Out</button>
        </div>
        <div class="button-group">
          <button type="button" class="primary-button" @click="createNewGroup">Create Group</button>
        </div>
      </div>
    </div>
    <ErrorMsg v-if="errorNotification" :msg="errorNotification" />
    <div>
      <p v-if="isLoading">Loading groups...</p>
      <div v-else-if="userGroups.length === 0">
        <p>No groups available.</p>
      </div>
      <div v-else class="groups-list-container">
        <div
          v-for="group in userGroups"
          :key="group.id"
          class="group-item-card"
          @click="navigateToGroup(group.id, group.name)"
        >
          <div class="group-image-container">
            <img
              v-if="group.conversationPhoto.String"
              :src="'data:image/png;base64,' + group.conversationPhoto.String"
              alt="Group Image"
              class="group-image"
            />
          </div>
          <div class="group-info">
            <h4>{{ group.name }}</h4>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header-container {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  padding-top: 1rem;
  padding-bottom: 0.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid #dee2e6;
}

.page-heading {
  font-size: 1.75rem;
  font-weight: 600;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
}

.button-group {
  position: relative;
  display: inline-flex;
  vertical-align: middle;
  margin-right: 0.5rem;
  margin-bottom: 0.5rem;
}

.secondary-button {
  padding: 0.375rem 0.75rem;
  background-color: transparent;
  color: #6c757d;
  border: 1px solid #6c757d;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.15s ease-in-out;
}

.secondary-button:hover {
  background-color: #6c757d;
  color: white;
}

.primary-button {
  padding: 0.375rem 0.75rem;
  background-color: transparent;
  color: #007bff;
  border: 1px solid #007bff;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.15s ease-in-out;
}

.primary-button:hover {
  background-color: #007bff;
  color: white;
}

.groups-list-container {
  display: flex;
  flex-direction: column;
}

.group-item-card {
  background-color: #f8f9fa;
  padding: 1rem;
  margin-bottom: 0.75rem;
  cursor: pointer;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: background-color 0.2s ease;
}

.group-item-card:hover {
  background-color: #e9ecef;
}

.group-image-container {
  flex-shrink: 0;
  width: 75px;
  height: 75px;
}

.group-image {
  width: 75px;
  height: 75px;
  object-fit: cover;
  border-radius: 50%;
}

.group-info h4 {
  margin-top: 0;
  margin-bottom: 0;
  color: #333;
}

@media (max-width: 600px) {
  .group-item-card {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-container {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .action-buttons {
    width: 100%;
    justify-content: space-between;
  }
}
</style>
