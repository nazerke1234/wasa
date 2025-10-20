<template>
  <div class="group-settings-panel">
    <div class="settings-header">
      <div class="group-image-wrapper">
        <img v-if="groupImage" :src="groupImage" alt="Group Image" class="group-image" />
      </div>
      <div class="group-details">
        <h1 class="group-title">{{ groupTitle }}</h1>
        <div class="name-update-section">
          <input
            v-model="updatedGroupTitle"
            placeholder="Enter new group title"
            maxlength="16"
            minlength="3"
          />
          <button
            @click="modifyGroupTitle"
            :disabled="!updatedGroupTitle || updatedGroupTitle === groupTitle"
          >
            Change Group Title
          </button>
        </div>
        <div class="image-update-section">
          <input type="file" @change="processImageSelection" accept="image/*" />
          <button @click="modifyGroupImage" :disabled="!newGroupImage">Change Group Image</button>
        </div>
        <div class="member-addition-section">
          <h3>Include New Members</h3>
          <form @submit.prevent="findUsers" class="user-search-form">
            <input
              v-model="searchQuery"
              placeholder="Find users by username"
              class="search-field"
            />
            <button type="submit" class="search-submit-btn">Find</button>
          </form>
          
          <div v-if="displaySearchResults" class="search-result-area">
            <div v-if="searchResults.length === 0" class="empty-results">
              No users found for "{{ previousSearchQuery }}"
            </div>
            <div v-else class="result-list">
              <div v-for="user in searchResults" :key="user.id" class="result-item">
                <span class="result-username">{{ user.name }}</span>
                <button
                  class="include-btn"
                  @click="includeUserInGroup(user.id)"
                  :disabled="isCurrentMember(user.id)"
                >
                  {{ isCurrentMember(user.id) ? 'Already in group' : 'Include in Group' }}
                </button>
              </div>
            </div>
          </div>
        </div>
        <div class="group-exit-section">
          <button class="exit-btn" @click="exitGroup">
            Exit Group
          </button>
        </div>
      </div>
    </div>
    <ErrorDisplay v-if="errorMessage" :msg="errorMessage" />
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorDisplay from "../components/ErrorMsg.vue";

export default {
  name: "GroupSettingsPanel",
  components: {
    ErrorDisplay,
  },
  data() {
    return {
      authToken: localStorage.getItem("token"),
      groupIdentifier: this.$route.params.uuid,
      groupTitle: localStorage.getItem("groupName"),
      groupImage: null,
      updatedGroupTitle: "", 
      newGroupImage: null, 
      errorMessage: null, 
      searchQuery: "",
      previousSearchQuery: "",
      searchResults: [],
      isLoading: false,
      displaySearchResults: false,
      groupMembers: [],
    };
  },
  methods: {
    async loadGroupDetails() {
      try {
        const userToken = localStorage.getItem("token");
        if (!userToken) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await axios.get(`/groups/${this.groupIdentifier}`, {
          headers: {
            Authorization: `Bearer ${userToken}`,
          },
        });
        const groupImageData = response.data.groupPhoto;
        this.groupImage = groupImageData ? `data:image/*;base64,${groupImageData}` : null;
        this.groupMembers = response.data.members;
      } catch (error) {
        console.error("Unable to load group details:", error);
        this.errorMessage = "Unable to load group information. Please attempt again later.";
      }
    },
    processImageSelection(event) {
      const selectedFile = event.target.files[0];
      if (selectedFile) {
        this.newGroupImage = selectedFile;
      }
    },
    async modifyGroupImage() {
      if (!this.newGroupImage) return;
      try {
        const imageData = new FormData();
        imageData.append("photo", this.newGroupImage);
        await axios.put(`/groups/${this.groupIdentifier}/photo`, imageData, {
          headers: {
            Authorization: `Bearer ${this.authToken}`,
          },
        });
        alert("Group image successfully updated!");
        this.newGroupImage = null;
        this.loadGroupDetails();
      } catch (error) {
        console.error("Unable to update group image:", error);
        this.errorMessage = "Unable to update group image. Please attempt again.";
      }
    },
    async modifyGroupTitle() {
      if (!this.updatedGroupTitle || this.updatedGroupTitle === this.groupTitle) return;
      try {
        await axios.put(
          `/groups/${this.groupIdentifier}/name`,
          { groupName: this.updatedGroupTitle },
          {
            headers: {
              Authorization: `Bearer ${this.authToken}`,
            },
          }
        );
        alert("Group title successfully updated!");
        localStorage.setItem("groupName", this.updatedGroupTitle);
        this.groupTitle = this.updatedGroupTitle;
        this.updatedGroupTitle = "";
      } catch (error) {
        console.error("Unable to update group title:", error);
        this.errorMessage = "Unable to update group title. Please attempt again.";
      }
    },
    async exitGroup() {
      if (!confirm('Are you certain you wish to exit this group?')) {
        return;
      }
      try {
        await axios.delete(`/groups/${this.groupIdentifier}`, {
          headers: {
            Authorization: `Bearer ${this.authToken}`,
          },
        });
        this.$router.push({ path: "/groups" });
      } catch (error) {
        console.error("Unable to exit group:", error);
        this.errorMessage = "Unable to exit group. Please attempt again.";
      }
    },  
    async findUsers() {
      if (!this.searchQuery.trim()) {
        this.errorMessage = "Please provide a search term";
        this.displaySearchResults = false;
        return;
      }

      this.isLoading = true;
      this.errorMessage = null;
      try {
        const response = await axios.get(`/search`, {
          params: { username: this.searchQuery }
        });
        this.searchResults = response.data;
        this.previousSearchQuery = this.searchQuery;
        this.displaySearchResults = true;
      } catch (error) {
        console.error("Search operation failed:", error);
        this.errorMessage = "Unable to find users. Please attempt again.";
      } finally {
        this.isLoading = false;
      }
    },

    isCurrentMember(userId) {
      return this.groupMembers.includes(userId);
    },

    async includeUserInGroup(userId) {
      if (this.isCurrentMember(userId)) return;
      try {
        await axios.post(`/groups/${this.groupIdentifier}`, 
        {  
          userId: userId,
        },
        {
            headers: {
            Authorization: `Bearer ${this.authToken}`,
            },}
          );
        this.groupMembers.push(userId);
        this.errorMessage = null;
      } catch (error) {
        console.error("Unable to include user:", error);
        this.errorMessage = "Unable to include user in group. Please attempt again.";
      }
  },
  },
  mounted() {
    this.loadGroupDetails();
  },
};
</script>

<style scoped>
.member-addition-section {
  margin-top: 2rem;
  padding: 1rem;
  border-top: 1px solid #eee;
}

.user-search-form {
  display: flex;
  gap: 1rem;
  margin: 1rem 0;
}

.search-field {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.search-submit-btn {
  padding: 0.5rem 1rem;
  background-color: #28a745;
  color: white;
}

.search-result-area {
  margin-top: 1rem;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.include-btn {
  padding: 0.25rem 0.75rem;
  background-color: #17a2b8;
  color: white;
}

.include-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.empty-results {
  color: #666;
  font-style: italic;
  padding: 1rem;
  text-align: center;
}

.group-settings-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.settings-header {
  display: flex;
  align-items: flex-start;
  gap: 20px;
  width: 100%;
  max-width: 800px;
}

.group-image-wrapper {
  flex: 0 0 auto;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid #ccc;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
}

.group-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image-placeholder {
  color: #aaa;
  font-size: 14px;
}

.group-details {
  flex: 1;
}

.group-title {
  margin: 0;
  font-size: 24px;
  font-weight: bold;
}

.name-update-section,
.image-update-section {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
}

input {
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}

button {
  padding: 8px 12px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

button:hover:not(:disabled) {
  background-color: #0056b3;
}
</style>
