<script>
export default {
  name: "GroupView",
  data() {
    return {
      username: "",
      errormsg: null,
      loading: false,
      groups: []
    };
  },
  methods: {
    async loadGroups() {
      this.errormsg = null;
      this.loading = true;
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await this.$axios.get("/groups", {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.groups = response.data || [];
      } catch (error) {
        console.error("Error loading groups:", error);
        this.errormsg = "Failed to load groups. Please try again.";
      } finally {
        this.loading = false;
      }
    },
    viewGroup(groupId, groupName) {
        localStorage.setItem("groupName", groupName);
        this.$router.push({
            path: `/groups/${groupId}`
        });
    },
    startChat() {
        this.$router.push({ path: "/search" });
    },
    logOut() {
        this.$router.push({ path: "/" });
    },
    newGroup() {
        this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    this.username = localStorage.getItem("name") || "Guest";
    this.loadGroups();
  }
};
</script>

<template>
  <div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2 conversation-header">All groups</h1>
      <div class="btn-toolbar mb-2 mb-md-0 button-group-container">
        
        <button type="button" class="btn-chat-primary" @click="startChat">
          Start Chat
        </button>

        <button type="button" class="btn-chat-secondary" @click="newGroup">
          New group
        </button>

        <button type="button" class="btn-chat-tertiary" @click="logOut">
          Log Out
        </button>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" />

    <div>
      <p v-if="loading">Loading...</p>
      <div v-else-if="groups.length === 0">
        <p>No groups found. Click "New group" to create one.</p>
      </div>
      <div v-else class="groups-container">
        <div
          v-for="group in groups"
          :key="group.id"
          class="group-block"
          @click="viewGroup(group.id, group.name)"
        >
          <div class="group-photo">
            <img
              v-if="group.conversationPhoto.String"
              :src="'data:image/png;base64,' + group.conversationPhoto.String"
              alt="Group Photo"
              class="group-picture"
            />
             <div v-else class="group-picture placeholder">
                 #
             </div>
          </div>
          <div class="group-details">
            <h4>{{ group.name }}</h4>
            <div class="members-list">
                 <span class="member-label">Members:</span>
                 <span v-if="group.members && group.members.length" class="member-names">
                     {{ group.members.join(', ') }}
                 </span>
                 <span v-else class="member-names none">Only you</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.conversation-header {
    font-size: 2rem;
    font-weight: 600;
    color: #2c3e50;
    margin: 0;
}

.button-group-container {
    display: flex;
    gap: 10px;
}

/* Start Chat */
.btn-chat-primary {
    background-color: #4a90e2;
    color: white;
    border: 1px solid #4a90e2;
    padding: 10px 18px;
    font-weight: 600;
    border-radius: 8px;
    transition: background-color 0.2s, box-shadow 0.2s;
    cursor: pointer;
}

.btn-chat-primary:hover {
    background-color: #3a74b6;
    box-shadow: 0 4px 10px rgba(74, 144, 226, 0.3);
}

/* New group */
.btn-chat-secondary {
    background-color: white;
    color: #4a90e2;
    border: 1px solid #4a90e2;
    padding: 10px 18px;
    font-weight: 600;
    border-radius: 8px;
    transition: background-color 0.2s, color 0.2s;
    cursor: pointer;
}

.btn-chat-secondary:hover {
    background-color: #e6f1fc;
    color: #3a74b6;
}

/* Log Out */
.btn-chat-tertiary {
    background-color: transparent;
    color: #7f8c8d;
    border: 1px solid #bdc3c7;
    padding: 10px 18px;
    font-weight: 500;
    border-radius: 8px;
    transition: background-color 0.2s, color 0.2s;
    cursor: pointer;
}

.btn-chat-tertiary:hover {
    background-color: #ecf0f1;
    color: #34495e;
    border-color: #34495e;
}
.groups-container {
  display: flex;
  flex-direction: column;
}

.group-block {
  
  background-color: white; 
  padding: 15px;
  margin-bottom: 12px; 
  cursor: pointer;
  border-radius: 10px; 
  display: flex;
  align-items: center;  
  gap: 15px;  
  border: 1px solid #ecf0f1; 
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05); 
  transition: all 0.2s ease;
}

.group-block:hover {
    background-color: #f7f9fc; 
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
}

.group-photo {
  flex-shrink: 0;  
  width: 50px; 
  height: 50px; 
}

.group-picture {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 50%;
  border: 2px solid #4a90e2; 
}

.group-picture.placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #e6f1fc;
    color: #4a90e2;
    font-size: 1.5rem;
    font-weight: 700;
}

.group-details h4 {
  margin-top: 0;
  margin-bottom: 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.members-list {
    margin-top: 4px;
    font-size: 0.85rem;
    color: #7f8c8d;
    display: flex;
    gap: 5px;
}

.member-label {
    font-weight: 600;
    color: #95a5a6;
}

.member-names {
    display: inline-block;
    max-width: 250px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.member-names.none {
    font-style: italic;
}

.group-members {
    margin: 4px 0 0 0;
    font-size: 0.85rem;
    color: #7f8c8d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 300px;
}

.empty-state {
    padding: 40px;
    text-align: center;
    color: #95a5a6;
    font-style: italic;
}



@media (max-width: 600px) {
  .group-block {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
