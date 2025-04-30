<template>
  <div class="chat-container">
    <div class="chat-header">
      <h2>Chat</h2>
    </div>

    <!-- Message List -->
    <div class="messages">
      <div
        v-for="message in messages"
        :key="message.ID"
        :class="{ 'my-message': message.isMine, 'other-message': !message.isMine }"
      >
        <div class="message-content">
          <span>{{ message.Content }}</span>

          <button v-if="message.isMine" @click="deleteMessage(message.ID)">
            <font-awesome-icon icon="trash" />
          </button>

          <button @click="message.showForwardInput = !message.showForwardInput">
            <font-awesome-icon icon="share" />
          </button>

          <!-- Add/remove comment button -->
          <button v-if="!message.showCommentInput" @click="toggleCommentInput(message)">
            <font-awesome-icon :icon="message.Comment ? 'times' : 'comment'" />
          </button>
        </div>

        <!-- Forward input field -->
        <div v-if="message.showForwardInput" class="forward-box">
          <input v-model="message.forwardChatName" type="text" placeholder="Chat name..." />
          <button @click="forwardMessage(message)">
            <font-awesome-icon icon="paper-plane" />
          </button>
        </div>
        
        <div v-if="message.showCommentInput" class="comment-box">
          <input v-model="message.newComment" type="text" placeholder="Enter comment..." />
          <button @click="addOrRemoveComment(message)">
            <font-awesome-icon icon="check" />
          </button>
        </div>

        <div v-if="message.Comment" class="comment">
          <strong>Comment:</strong> {{ message.Comment }}
        </div>

      </div>
    </div>

    <!-- Message input form -->
    <div class="message-input">
      <input 
        v-model="newMessage"
        type="text"
        placeholder="Enter message..."
        @keyup.enter="sendMessage"
      />
      <button @click="sendMessage">
        <font-awesome-icon icon="paper-plane" />
      </button>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { library } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import { faTrash, faShare, faComment, faTimes, faCheck, faPaperPlane } from "@fortawesome/free-solid-svg-icons";

library.add(faTrash, faShare, faComment, faTimes, faCheck, faPaperPlane);

export default {
  components: { FontAwesomeIcon },
  setup() {
    const route = useRoute();
    const chatId = route.params.chatId;
    const messages = ref([]);
    const newMessage = ref("");

    const token = localStorage.getItem("access_token")?.trim();
    const userId = JSON.parse(localStorage.getItem("user_id")) || null;

    if (!token) {
      console.error("Token missing! User is not authorized.");
    }

    const axiosConfig = {
      headers: { Authorization: `Bearer ${token}` },
    };

    const fetchMessages = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/api/chats/${chatId}`, axiosConfig);
        messages.value = response.data.map((msg) => ({
          ...msg,
          isMine: msg.SenderID === userId,
          showForwardInput: false,
          forwardUsername: "",
        }));
      } catch (error) {
        console.error("Error fetching messages:", error.response?.data || error.message);
      }
    };

    const sendMessage = async () => {
      if (!newMessage.value.trim()) return;

      try {
        const messageData = {
          chat_id: Number(chatId),
          content: newMessage.value.trim(),
        };

        const response = await axios.post("http://localhost:8080/api/messages", messageData, axiosConfig);

        messages.value.push({
          ID: response.data.ID,
          Content: newMessage.value,
          SenderID: userId,
          isMine: true,
          showForwardInput: false,
          forwardUsername: "",
        });

        newMessage.value = "";
        fetchMessages();
      } catch (error) {
        console.error("Error sending message:", error.response?.data || error.message);
      }
    };

    const forwardMessage = async (message) => {
      if (!message.forwardChatName.trim()) {
        alert("Please enter the chat name!");
        return;
      }

      try {
        await axios.post(
          "http://localhost:8080/api/messages/forward",
          { 
            message_id: message.ID, 
            chat_name: message.forwardChatName, 
          },
          axiosConfig
        );

        alert("Message forwarded!");
        message.forwardChatName = ""; // Clear the input field
      } catch (error) {
        alert("Chat doesnt exist!");
        console.error("Error forwarding message:", error.response?.data || error.message);
      }
    };

    const toggleCommentInput = (message) => {
      message.showCommentInput = !message.showCommentInput;
    };

    const addOrRemoveComment = async (message) => {
      if (message.Comment) {
        await uncommentMessage(message.ID);
      } else {
        await commentMessage(message.ID, message.newComment);
      }
      message.showCommentInput = false;
    };

    const commentMessage = async (messageId, commentText) => {
      if (!commentText?.trim()) return;

      try {
        await axios.put(
          `http://localhost:8080/api/messages/${messageId}/comment`,
          { comment: commentText },
          axiosConfig
        );

        const msg = messages.value.find((m) => m.ID === messageId);
        if (msg) msg.Comment = commentText;
        fetchMessages();
      } catch (error) {
        console.error("Error adding comment:", error.response?.data || error.message);
      }
    };

    const uncommentMessage = async (messageId) => {
      try {
        await axios.delete(`http://localhost:8080/api/messages/${messageId}/comment`, axiosConfig);

        const msg = messages.value.find((m) => m.ID === messageId);
        if (msg) msg.Comment = null;
        fetchMessages();
      } catch (error) {
        console.error("Error deleting comment:", error.response?.data || error.message);
      }
    };

    const deleteMessage = async (messageId) => {
      try {
        await axios.delete(`http://localhost:8080/api/messages/${messageId}`, axiosConfig);

        messages.value = messages.value.filter((msg) => msg.ID !== messageId);
        fetchMessages();
      } catch (error) {
        console.error("Error deleting message:", error.response?.data || error.message);
      }
    };

    onMounted(() => {
      fetchMessages();
    });

    return {
      messages,
      newMessage,
      sendMessage,
      deleteMessage,
      forwardMessage,
      commentMessage,
      uncommentMessage,
      toggleCommentInput,
      addOrRemoveComment,
    };
  },
};
</script>


<style scoped>
  .chat-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #1e2d2b; 
  color: white;
}

.chat-header {
  padding: 10px;
  background: #223432;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
}

.messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
}

.my-message, .other-message {
  padding: 8px 12px;
  border-radius: 12px;
  max-width: 60%;
  margin-bottom: 10px;
  word-wrap: break-word;
}

.my-message {
  background: #6a5acd; 
  align-self: flex-end;
  color: white;
}

.other-message {
  background: #444; 
  align-self: flex-start;
  color: white;
}

.comment {
  font-size: 14px;
  color: #bbb;
  margin-top: 5px;
}

.message-input {
  display: flex;
  padding: 10px;
  background: #223432;
  border-top: 1px solid #333;
}

.message-input input {
  flex: 1;
  padding: 10px;
  border: none;
  border-radius: 20px;
  background: #2e3f3c;
  color: white;
  outline: none;
}

.message-input button {
  padding: 10px 15px;
  margin-left: 10px;
  border: none;
  background: #6a5acd;
  color: white;
  border-radius: 10px;
  cursor: pointer;
}

.forward-box {
  display: flex;
  gap: 5px;
  margin-top: 5px;
}

.forward-box input {
  flex: 1;
  padding: 6px;
  border-radius: 5px;
  background: #2e3f3c;
  color: white;
  border: none;
}

.forward-box button {
  background: #6a5acd;
  color: white;
  border: none;
  padding: 6px;
  border-radius: 5px;
  cursor: pointer;
}

.search-users, .profile-btn {
  position: absolute;
  top: 10px;
  font-size: 14px;
  color: white;
  padding: 5px 10px;
  border-radius: 5px;
}

.search-users {
  left: 10px;
  background: #3b5249;
}

.profile-btn {
  right: 10px;
  background: #3b5249;
}

.comment-box {
  display: flex;
  gap: 5px;
  margin-top: 5px;
}

.comment-box input {
  flex: 1;
  padding: 6px;
  border-radius: 5px;
  background: #2e3f3c;
  color: white;
  border: none;
}

.comment-box button {
  background: #6a5acd;
  color: white;
  border: none;
  padding: 6px;
  border-radius: 5px;
  cursor: pointer;
}

</style>
