<template>
  <div class="group-settings-panel">
    
    <LoadingSpinner v-if="isLoading" />
    
    <div class="settings-header">
      
      <div class="group-image-wrapper">
        <img v-if="groupImage" :src="groupImage" alt="Group Image" class="group-image" />
        <div v-else class="group-image no-image-placeholder">
           No Photo
        </div>
      </div>

      <div class="group-details">
        <h1 class="group-title">{{ groupName }}</h1>
        
        <div class="update-section name-update-section">
          <input
            v-model="updatedGroupTitle"
            placeholder="Enter new group title (3-16 chars)"
            maxlength="16"
            minlength="3"
            class="input-field"
          />
          <button
            class="btn-primary"
            @click="modifyGroupTitle"
            :disabled="!updatedGroupTitle || updatedGroupTitle === groupName"
          >
            Change Name
          </button>
        </div>
        
        <div class="update-section image-update-section">
          <input type="file" @change="processImageSelection" accept="image/jpeg, image/png" class="file-input" />
          <button class="btn-primary" @click="modifyGroupImage" :disabled="!newGroupImage">
            Change Image
          </button>
        </div>
      </div>
    </div>
    
    <ErrorMsg v-if="errorMessage" :msg="errorMessage" />
    <SuccessMsg v-if="successMessage" :msg="successMessage" />

    <div class="member-addition-section">
      <h3 class="section-title">Include New Members</h3>
      
      <form @submit.prevent="findUsers" class="user-search-form">
        <input
          v-model="searchQuery"
          placeholder="Find users by username"
          class="search-field input-field"
        />
        <button type="submit" class="btn-secondary search-submit-btn" :disabled="isLoading">
          <LoadingSpinner v-if="isLoading" small />
          <span v-else>Find</span>
        </button>
      </form>
      
      <div v-if="displaySearchResults" class="search-result-area">
        <div v-if="searchResults.length === 0" class="empty-results">
          No users found for "{{ previousSearchQuery }}"
        </div>
        <div v-else class="result-list">
          <div v-for="user in searchResults" :key="user.id" class="result-item">
            <span class="result-username">{{ user.name }}</span>
            <button
              class="btn-tertiary include-btn"
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
      <button class="btn-danger" @click="exitGroup">
        Exit Group
      </button>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";

export default {
  name: "EditGroupView",

  data() {
    return {
      token: localStorage.getItem("token"),
      groupId: this.$route.params.uuid,
      groupName: localStorage.getItem("groupName"),
      groupImage: null,
      updatedGroupTitle: "",
      newGroupImage: null,
      errorMessage: null,
      successMessage: null, 
      isLoading: false, 
      searchQuery: "",
      previousSearchQuery: "",
      searchResults: [],
      displaySearchResults: false,
      groupMembers: [],
    };
  },
  methods: {
    setMessages(type, message) {
        if (type === 'success') {
            this.successMessage = message;
            this.errorMessage = null;
        } else if (type === 'error') {
            this.errorMessage = message;
            this.successMessage = null;
        } else {
            this.successMessage = null;
            this.errorMessage = null;
        }
        
        if (type === 'success') {
            setTimeout(() => {
                this.successMessage = null;
            }, 3000); 
        }
    },
    async loadGroupDetails() {
      this.setMessages('clear');
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await axios.get(`/groups/${this.groupId}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const photoBase64 = response.data.groupPhoto;
        
        this.groupImage = photoBase64 ? `data:image/png;base64,${photoBase64}` : null;
        this.groupMembers = response.data.members;
        this.groupName = response.data.name; 
        localStorage.setItem("groupName", this.groupName);
      } catch (error) {
        console.error("Unable to load group details:", error);
        this.setMessages('error', "Unable to load group information. Please attempt again later.");
      }
    },
    processImageSelection(event) {
      this.setMessages('clear');
      const selectedFile = event.target.files[0];
      if (selectedFile) {
        this.newGroupImage = selectedFile;
      }
    },
    async modifyGroupImage() {
      if (!this.newGroupImage) return;
      this.setMessages('clear');
      this.isLoading = true;
      try {
        const imageData = new FormData();
        imageData.append("photo", this.newGroupImage);
        await axios.put(`/groups/${this.groupId}/photo`, imageData, {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },
        });
        this.setMessages('success', "Group image successfully updated!");
        this.newGroupImage = null;
        this.loadGroupDetails();
      
      } catch (error) {
        console.error("Unable to update group image:", error);
        this.setMessages('error', "Unable to update group image. Please attempt again.");
      } finally {
        this.isLoading = false;
      }
    },
    async modifyGroupTitle() {
      if (!this.updatedGroupTitle || this.updatedGroupTitle === this.groupName) return;
      this.setMessages('clear');
      this.isLoading = true;
      try {
        await axios.put(
          `/groups/${this.groupId}/name`,
          { groupName: this.updatedGroupTitle },
          {
            headers: {
              Authorization: `Bearer ${this.token}`,
            },
          }
        );
        this.setMessages('success', "Group name successfully updated!");
        localStorage.setItem("groupName", this.updatedGroupTitle);
        this.groupName = this.updatedGroupTitle;
        this.updatedGroupTitle = "";
      } catch (error) {
        console.error("Unable to update group title:", error);
        this.setMessages('error', "Unable to update group title. Please attempt again.");
      } finally {
        this.isLoading = false;
      }
    },
    async exitGroup() {
      if (!confirm('Are you sure you wish to exit this group?')) 
      {
        return;
      }
      try {
        await axios.delete(`/groups/${this.groupId}`, {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },
        });
        this.$router.push({ path: "/groups" });
      } catch (error) {
        console.error("Unable to exit group:", error);
        this.setMessages('error', "Unable to exit group. Please attempt again.");
      } finally {
        this.isLoading = false;
      }
    },  
    async findUsers() {
      if (!this.searchQuery.trim()) {
        this.setMessages('error', "Please provide a search term.");
        this.displaySearchResults = false;
        return;
      }
      this.setMessages('clear');
      this.isLoading = true;
      try {
        const response = await axios.get(`/search`, {
          params: { username: this.searchQuery }
        });
        this.searchResults = response.data;
        this.previousSearchQuery = this.searchQuery;
        this.displaySearchResults = true;
      } catch (error) {
        console.error("Search operation failed:", error);
        this.setMessages('error', "Unable to find users. Please attempt again.");
      } finally {
        this.isLoading = false;
      }
    },

    isCurrentMember(userId) {
      return this.groupMembers.includes(userId);
    },

    async includeUserInGroup(userId) {
      if (this.isCurrentMember(userId)) return;
      this.setMessages('clear');
      try {
        await axios.post(`/groups/${this.groupId}`, 
        {  
          userId: userId,
        },
        {
            headers: {
            Authorization: `Bearer ${this.token}`,
            },}
          );
        this.groupMembers.push(userId);
        this.setMessages('success', "User successfully added to the group!");
        this.displaySearchResults = false;
        this.searchResults = [];
    
      } catch (error) {
        console.error("Unable to include user:", error);
        this.setMessages('error', "Unable to include user in group. Please attempt again.");
      }
  },
  },
  mounted() {
    this.updatedGroupTitle = this.groupName; 
    this.loadGroupDetails();
  },
};
</script>

<style scoped>

.group-settings-panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px 20px;
  background-color: #f7f9fb;
  border-radius: 12px;
  max-width: 900px;
  margin: 20px auto;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.settings-header {
  display: flex;
  align-items: flex-start;
  gap: 30px;
  width: 100%;
  padding-bottom: 20px;
  border-bottom: 1px solid #e0e0e0; /* Fixed: --color-border-light */
  margin-bottom: 20px;
}

.group-details {
  flex: 1;
}

.group-title {
  margin: 0 0 15px 0;
  font-size: 2.2rem;
  font-weight: 700;
  color: #34495e; /* Fixed: --color-text-dark */
}


.group-image-wrapper {
  flex: 0 0 auto;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  border: 4px solid #4a90e2; /* Fixed: --brand-color-primary */
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ffffff;
}

.group-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.no-image-placeholder {
  color: #aaa;
  font-size: 1rem;
  text-align: center;
  font-weight: 600;
}


.update-section {
  margin-top: 15px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.input-field {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #ccc;
  border-radius: 6px;
  transition: border-color 0.2s;
  font-size: 1rem;
}
.input-field:focus {
  border-color: #4a90e2; /* Fixed: --brand-color-primary */
  outline: none;
  box-shadow: 0 0 0 2px rgba(74, 144, 226, 0.2);
}

.file-input {
  flex: 1;
  border: none;
  padding: 0;
  font-size: 0.9rem;
}


button {
  padding: 10px 18px;
  font-weight: 600;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s, opacity 0.2s;
  white-space: nowrap;
}

.btn-primary {
  background-color: #4a90e2; /* Fixed: --brand-color-primary */
  color: white;
}
.btn-primary:hover:not(:disabled) {
  background-color: #3a74b6;
}

.btn-secondary {
  background-color: #28a745; /* Fixed: --brand-color-secondary */
  color: white;
}
.btn-secondary:hover:not(:disabled) {
  background-color: #1f8b36;
}

.btn-tertiary {
  background-color: #f0f0f0;
  color: #4a90e2; /* Fixed: --brand-color-primary */
  border: 1px solid #4a90e2; /* Fixed: --brand-color-primary */
  padding: 8px 14px;
  font-size: 0.9rem;
}
.btn-tertiary:hover:not(:disabled) {
  background-color: #e6f1fc;
}

.btn-danger {
  background-color: #dc3545; /* Fixed: --brand-color-danger */
  color: white;
  margin-top: 30px;
}
.btn-danger:hover:not(:disabled) {
  background-color: #c82333;
}

button:disabled {
  background-color: #ccc;
  opacity: 0.7;
  cursor: not-allowed;
}


.member-addition-section {
  width: 100%;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px dashed #e0e0e0; /* Fixed: --color-border-light */
}

.section-title {
  color: #34495e; /* Fixed: --color-text-dark */
  margin-top: 0;
}

.user-search-form {
  display: flex;
  gap: 15px;
  margin: 15px 0;
}

.search-field {
  flex: 1;
}

.search-result-area {
  margin-top: 20px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background-color: #fff;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}
.result-item:last-child {
  border-bottom: none;
}

.result-username {
  font-weight: 500;
  color: #34495e; /* Fixed: --color-text-dark */
}

.empty-results {
  color: #95a5a6;
  font-style: italic;
  padding: 10px;
  text-align: center;
}


@media (max-width: 600px) {
  .settings-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  .group-image-wrapper {
    margin-bottom: 20px;
  }
  .update-section,
  .user-search-form {
    flex-direction: column;
    gap: 10px;
  }
  .input-field,
  button,
  .file-input {
    width: 100%;
  }
}
</style>

