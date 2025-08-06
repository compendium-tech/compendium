<template>
  <StandardLayout>
    <div class="min-h-screen bg-gradient-to-b from-gray-100 to-gray-200 py-48">
      <!-- Header -->
      <header class="text-center mb-12">
        <h1 class="text-5xl font-extrabold text-gray-900 mb-3">College Search</h1>
        <p class="text-xl text-gray-600">Discover your dream college with our AI-powered search!</p>
      </header>

      <!-- Search Form -->
      <div class="max-w-4xl mx-auto bg-white rounded-2xl shadow-2xl p-8 mb-12">
        <h2 class="text-3xl font-semibold text-gray-800 mb-8">Search Parameters</h2>
        <form @submit.prevent="searchColleges" class="flex-col space-y-6">
          <div class="relative">
            <label for="semanticSearchText" class="block text-lg font-medium text-gray-700">AI-Powered Search</label>
            <p class="text-sm text-gray-500 mb-3">Let our AI find colleges tailored to your dreams!</p>
            <BaseInput id="semanticSearchText" v-model="form.semanticSearchText" type="text"
              placeholder="Search by college name, program, or vibe" />
            <span class="rainbow-text absolute -top-8 left-0 text-2xl font-bold animate-pulse">
              AI Magic Search ✨
            </span>
          </div>
          <div>
            <label for="stateOrCountry" class="block text-lg font-medium text-gray-700">State or Country</label>
            <p class="text-sm text-gray-500 mb-3">e.g. "California", "USA"</p>
            <BaseInput id="stateOrCountry" v-model="form.stateOrCountry" type="text" placeholder="Filter by location" />
          </div>
          <div class="flex justify-end space-x-4">
            <BaseButton type="button" variant="secondary" size="md" :disabled="!searchHistory.length"
              @click="goToPreviousSearch" hover-effect="translate">
              Previous Search
            </BaseButton>
            <BaseButton type="submit" variant="primary" size="md" :disabled="isLoading" hover-effect="translate"
              class="animate-bounce">
              <span v-if="isLoading">Searching...</span>
              <span v-else>Search Colleges</span>
            </BaseButton>
          </div>
        </form>
      </div>

      <!-- Search Results -->
      <div class="max-w-4xl mx-auto bg-white rounded-2xl shadow-2xl p-8">
        <h2 class="text-3xl font-semibold text-gray-800 mb-8">Results</h2>

        <!-- Loading and Error States -->
        <div v-if="isLoading" class="flex justify-center items-center py-12">
          <svg class="animate-spin -ml-1 mr-3 h-12 w-12 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none"
            viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
            </path>
          </svg>
        </div>
        <p v-else-if="error" class="text-red-600 text-center py-12">{{ error }}</p>
        <div v-else-if="results.length > 0" class="space-y-8">
          <div v-for="college in results" :key="college.name"
            class="p-6 bg-gray-50 rounded-lg border border-gray-200 hover:shadow-xl transition-shadow duration-300 cursor-pointer"
            @click="openPopup(college)">
            <h3 class="text-2xl font-bold text-gray-800">{{ college.name }}</h3>
            <p class="text-sm text-gray-500 mb-3">{{ college.city }}, {{ college.stateOrCountry }}</p>
            <p class="text-gray-700 line-clamp-3">{{ college.description }}</p>
          </div>
        </div>
        <p v-else class="text-gray-500 text-center py-12">No colleges found. Try a different search query!</p>
      </div>

      <!-- Popup for College Details -->
      <div v-if="selectedCollege" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
        @click="closePopup">
        <div
          class="bg-white rounded-2xl p-8 max-w-3xl w-full max-h-[80vh] overflow-y-auto transform transition-all duration-500 scale-100"
          @click.stop>
          <h2 class="text-3xl font-bold text-gray-800 mb-4">{{ selectedCollege.name }}</h2>
          <p class="text-lg text-gray-600 mb-4">{{ selectedCollege.city }}, {{ selectedCollege.stateOrCountry }}</p>
          <div class="prose prose-sm" v-html="markedDescription"></div>
          <BaseButton class="mt-6" variant="primary" size="md" @click="closePopup" hover-effect="translate">
            Close
          </BaseButton>
        </div>
      </div>
    </div>
  </StandardLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from "vue";
import { Icon } from "@iconify/vue";
import BaseInput from "../ui/BaseInput.vue";
import BaseButton from "../ui/BaseButton.vue";
import StandardLayout from "../layout/StandardLayout.vue";
import BaseTransitioningText from "../ui/BaseTransitioningText.vue";
import { marked } from "marked";

// Mock API Client and Service
interface SearchCollegesRequest {
  pageIndex?: number;
  stateOrCountry?: string;
  semanticSearchText?: string;
}

interface CollegeResponse {
  name: string;
  city: string;
  stateOrCountry: string;
  description: string;
}

const mockColleges: CollegeResponse[] = [
  {
    name: "University of California, Berkeley",
    city: "Berkeley",
    stateOrCountry: "California",
    description: `
# University of California, Berkeley

## Overview
The University of California, Berkeley, fondly known as UC Berkeley or Cal, is a public research university established in 1868. As one of the flagship campuses of the University of California system, it has built a global reputation for academic excellence, groundbreaking research, and a vibrant campus culture that fosters innovation and social change.

## Academic Excellence
UC Berkeley offers over 350 degree programs across 14 colleges and schools, including highly ranked programs in engineering, computer science, environmental science, and social sciences. The university is home to the prestigious Haas School of Business, the College of Engineering, and the School of Public Health, among others. Its faculty includes numerous Nobel laureates, Pulitzer Prize winners, and MacArthur Fellows, ensuring students receive world-class instruction.

## Research Opportunities
Berkeley is a leader in research, with institutes like the Lawrence Berkeley National Laboratory and the Space Sciences Laboratory driving advancements in fields like physics, biology, and artificial intelligence. Students have unparalleled opportunities to engage in cutting-edge research, often working alongside faculty on projects that address global challenges such as climate change and public health.

## Campus Life
The campus, nestled in the scenic San Francisco Bay Area, is a hub of activity. With over 1,000 student organizations, including cultural clubs, professional societies, and activist groups, students can explore diverse interests. The iconic Sather Gate, Sproul Plaza, and the Campanile (Sather Tower) are central landmarks that define the Berkeley experience. The university's spirited athletics program, with its mascot Oski the Bear, competes in the NCAA Division I Pac-12 Conference.

## Diversity and Community
UC Berkeley prides itself on its diverse student body, with students from all 50 states and over 100 countries. The university is committed to equity and inclusion, offering resources like the Multicultural Community Center and programs to support first-generation and underrepresented students. Berkeley's progressive ethos is reflected in its history of activism, from the Free Speech Movement of the 1960s to modern-day advocacy for social justice.

## Location and Lifestyle
Located in Berkeley, California, the campus is just a short distance from San Francisco, offering students access to a vibrant urban environment with world-class dining, cultural attractions, and tech innovation hubs like Silicon Valley. The surrounding area features hiking trails, beaches, and a temperate climate, making it an ideal place for outdoor enthusiasts.

## Why Choose UC Berkeley?
UC Berkeley is more than a university—it's a community of thinkers, innovators, and change-makers. Whether you're passionate about advancing technology, tackling environmental challenges, or shaping public policy, Berkeley provides the resources, network, and inspiration to help you succeed. Its alumni include leaders in industry, government, and academia, making it a launchpad for impactful careers.
    `,
  },
  {
    name: "Stanford University",
    city: "Palo Alto",
    stateOrCountry: "California",
    description: `
# Stanford University

## Overview
Founded in 1885 by Leland and Jane Stanford, Stanford University is a private research institution located in the heart of Silicon Valley. Renowned for its entrepreneurial spirit, academic rigor, and contributions to technology and innovation, Stanford is a global leader in higher education.

## Academic Excellence
Stanford offers a wide range of undergraduate and graduate programs through its seven schools, including the School of Engineering, School of Humanities and Sciences, and the Graduate School of Business. Its interdisciplinary approach encourages students to explore diverse fields, from computer science to creative writing. The university is particularly celebrated for its programs in artificial intelligence, entrepreneurship, and biomedical sciences.

## Research and Innovation
Stanford is a powerhouse of innovation, with faculty and students contributing to advancements in technology, medicine, and sustainability. The university's proximity to Silicon Valley fosters partnerships with tech giants like Google, Apple, and Tesla, providing students with unique opportunities for internships and collaboration. Stanford's research centers, such as the Stanford Artificial Intelligence Laboratory and the Hoover Institution, drive global impact.

## Campus Life
The sprawling 8,180-acre campus is one of the largest in the United States, featuring stunning Spanish-style architecture, palm-lined pathways, and the iconic Hoover Tower. Students can join over 600 student organizations, participate in Division I athletics as part of the Pac-12 Conference, or engage in community service through the Haas Center for Public Service. The campus is also home to world-class facilities like the Stanford Arts District and the Anderson Collection.

## Diversity and Inclusion
Stanford is committed to fostering a diverse and inclusive community, with initiatives like the Diversity and First-Gen Office supporting students from varied backgrounds. The university hosts cultural events, such as the Asian American Activities Center’s annual Lunar New Year celebration, and provides resources for underrepresented groups to thrive academically and socially.

## Location and Lifestyle
Situated in Palo Alto, Stanford benefits from its proximity to Silicon Valley’s tech ecosystem and the cultural vibrancy of the San Francisco Bay Area. Students enjoy a Mediterranean climate, access to outdoor activities like hiking in the Santa Cruz Mountains, and a dynamic food scene. The campus’s bike-friendly infrastructure makes it easy to explore both the university and surrounding areas.

## Why Choose Stanford?
Stanford is a place where ideas are born and transformed into reality. Its emphasis on innovation, interdisciplinary learning, and global impact makes it an ideal choice for students who aspire to lead in their fields. With a network of accomplished alumni, including founders of companies like Cisco and Nike, Stanford equips students to shape the future.
    `,
  },
  {
    name: "Massachusetts Institute of Technology",
    city: "Cambridge",
    stateOrCountry: "Massachusetts",
    description: `
# Massachusetts Institute of Technology

## Overview
The Massachusetts Institute of Technology (MIT), founded in 1861, is a private research university in Cambridge, Massachusetts, known for its leadership in science, technology, engineering, and mathematics (STEM). MIT’s mission to advance knowledge and educate students in service of humanity has made it a global leader in innovation.

## Academic Excellence
MIT offers rigorous programs across five schools, including the School of Engineering, School of Science, and the MIT Sloan School of Management. Its hands-on approach to learning, exemplified by the motto “Mens et Manus” (Mind and Hand), emphasizes practical problem-solving. Programs in computer science, artificial intelligence, and physics are among the best in the world.

## Research and Innovation
MIT is at the forefront of technological advancement, with research centers like the MIT Media Lab, Lincoln Laboratory, and the Broad Institute pushing boundaries in AI, robotics, and genomics. Students have opportunities to work on real-world problems, from developing sustainable energy solutions to advancing quantum computing.

## Campus Life
MIT’s campus along the Charles River is a vibrant community of innovators. Students can join over 450 student organizations, participate in hackathons, or compete in Division III athletics. The campus’s unique traditions, like the MIT Mystery Hunt and the annual “Ring Premiere” for the iconic Brass Rat class ring, foster a strong sense of community.

## Diversity and Inclusion
MIT is dedicated to creating an inclusive environment, with programs like the Office of Minority Education and the Women’s and Gender Studies program supporting diverse students. The university hosts events celebrating cultural heritage and provides resources for first-generation and low-income students.

## Location and Lifestyle
Located in Cambridge, MIT is minutes from Boston, offering students access to a rich cultural and intellectual hub. The area is known for its historic charm, vibrant arts scene, and proximity to other top universities like Harvard. Students can enjoy kayaking on the Charles River, exploring Boston’s Freedom Trail, or attending concerts at nearby venues.

## Why Choose MIT?
MIT is a place for those who want to tackle the world’s toughest challenges through science and technology. Its collaborative culture, cutting-edge research, and emphasis on innovation make it a top choice for aspiring engineers, scientists, and entrepreneurs. MIT alumni, including Nobel laureates and tech pioneers, continue to shape the world.
    `,
  },
  {
    name: "Harvard University",
    city: "Cambridge",
    stateOrCountry: "Massachusetts",
    description: `
# Harvard University

## Overview
Founded in 1636, Harvard University is the oldest institution of higher education in the United States and a global leader in academic excellence. As a private Ivy League university, Harvard is renowned for its rigorous academics, distinguished faculty, and influential alumni network.

## Academic Excellence
Harvard offers a broad range of programs through its 13 schools, including Harvard College, Harvard Law School, and Harvard Medical School. Its liberal arts curriculum encourages exploration across disciplines, from history and literature to data science and public policy. The university’s faculty includes world-renowned scholars who mentor students in small, discussion-based classes.

## Research and Innovation
Harvard is a hub of intellectual discovery, with research centers like the Harvard Kennedy School’s Belfer Center for Science and International Affairs and the Wyss Institute for Biologically Inspired Engineering. Students can engage in research through programs like the Harvard College Research Program, tackling issues from climate change to global health.

## Campus Life
Harvard’s historic campus in Cambridge is home to iconic landmarks like Widener Library and Memorial Hall. With over 400 student organizations, including the Harvard Crimson newspaper and the Harvard Radcliffe Orchestra, students have endless opportunities to pursue their passions. Harvard competes in NCAA Division I athletics as part of the Ivy League.

## Diversity and Inclusion
Harvard is committed to diversity, with initiatives like the Harvard Foundation for Intercultural and Race Relations fostering an inclusive community. The university supports students from all backgrounds through financial aid programs, ensuring accessibility for low-income and first-generation students.

## Location and Lifestyle
Located in Cambridge, Harvard offers students access to Boston’s vibrant cultural and intellectual scene. The area is rich with history, museums, and restaurants, while nearby green spaces like the Boston Common provide opportunities for relaxation. The proximity to other universities fosters a collaborative academic environment.

## Why Choose Harvard?
Harvard is a place where tradition meets innovation. Its unparalleled resources, global network, and commitment to shaping leaders make it a top choice for students who want to make a difference. Harvard alumni, including U.S. presidents and Nobel laureates, exemplify the university’s legacy of impact.
    `,
  },
  {
    name: "University of Washington",
    city: "Seattle",
    stateOrCountry: "Washington",
    description: `
# University of Washington

## Overview
Founded in 1861, the University of Washington (UW) is a public research university in Seattle, known for its excellence in computer science, medicine, and environmental sciences. As one of the top public universities in the U.S., UW combines academic rigor with a commitment to public service.

## Academic Excellence
UW offers over 180 majors across 16 colleges and schools, including the highly ranked Paul G. Allen School of Computer Science & Engineering and the School of Medicine. Its interdisciplinary programs, such as the Environmental Science and Resource Management program, prepare students for global challenges.

## Research and Innovation
UW is a leader in research, with centers like the Institute for Health Metrics and Evaluation and the Clean Energy Institute driving advancements in health and sustainability. Students can participate in research through initiatives like the Undergraduate Research Program, working on projects that address real-world issues.

## Campus Life
The UW campus is renowned for its stunning beauty, with views of Mount Rainier and cherry blossoms in the Quad during spring. Students can join over 800 student organizations, from cultural clubs to outdoor adventure groups. The Huskies compete in NCAA Division I athletics as part of the Pac-12 Conference.

## Diversity and Inclusion
UW is committed to fostering a diverse community, with programs like the Office of Minority Affairs and Diversity supporting underrepresented students. The university hosts cultural events and provides resources to ensure all students feel included and empowered.

## Location and Lifestyle
Located in Seattle, UW offers students access to a vibrant city known for its tech industry, coffee culture, and outdoor recreation. Students can explore Pike Place Market, hike in the Cascade Mountains, or attend concerts at venues like the Showbox. The city’s mild climate is ideal for year-round outdoor activities.

## Why Choose UW?
The University of Washington is a place where students can pursue academic excellence while making a positive impact. Its strong ties to the tech industry, commitment to research, and vibrant campus life make it an ideal choice for students who want to innovate and lead.
    `,
  },
];

const collegeService = {
  searchColleges: async (request: SearchCollegesRequest): Promise<CollegeResponse[]> => {
    console.log("Searching with request:", request);
    return new Promise((resolve) => {
      setTimeout(() => {
        const searchText = request.semanticSearchText?.toLowerCase() || "";
        const stateOrCountry = request.stateOrCountry?.toLowerCase() || "";

        const filteredColleges = mockColleges.filter((college) => {
          const nameMatch = college.name.toLowerCase().includes(searchText);
          const descriptionMatch = college.description.toLowerCase().includes(searchText);
          const stateMatch = college.stateOrCountry.toLowerCase().includes(stateOrCountry);

          return (nameMatch || descriptionMatch) && stateMatch;
        });

        resolve(filteredColleges);
      }, 1000); // Simulate network latency
    });
  },
};

// Main application logic
const form = reactive<SearchCollegesRequest>({
  semanticSearchText: "",
  stateOrCountry: "",
});

const results = ref<CollegeResponse[]>([]);
const isLoading = ref(false);
const error = ref<string | null>(null);
const selectedCollege = ref<CollegeResponse | null>(null);
const searchHistory = ref<SearchCollegesRequest[]>([]);

const markedDescription = computed(() => {
  return selectedCollege.value ? marked(selectedCollege.value.description) : "";
});

const searchColleges = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    const response = await collegeService.searchColleges(form);
    results.value = response;
    searchHistory.value.push({ ...form });
  } catch (err) {
    console.error("API Error:", err);
    error.value = "Failed to fetch college data. Please try again.";
    results.value = [];
  } finally {
    isLoading.value = false;
  }
};

const openPopup = (college: CollegeResponse) => {
  selectedCollege.value = college;
};

const closePopup = () => {
  selectedCollege.value = null;
};

const goToPreviousSearch = () => {
  if (searchHistory.value.length > 0) {
    const lastSearch = searchHistory.value.pop();
    if (lastSearch) {
      form.semanticSearchText = lastSearch.semanticSearchText || "";
      form.stateOrCountry = lastSearch.stateOrCountry || "";
      searchColleges();
    }
  }
};
</script>

<style scoped>
.rainbow-text {
  background: linear-gradient(90deg,
      #ff0000,
      #ff9900,
      #33cc33,
      #00ccff,
      #cc00ff,
      #ff0000);
  background-size: 200%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: rainbow 8s ease infinite;
}

@keyframes rainbow {
  0% {
    background-position: 0%;
  }

  50% {
    background-position: 200%;
  }

  100% {
    background-position: 0%;
  }
}

.fade-height-enter-active,
.fade-height-leave-active {
  transition: all 0.3s ease-in-out;
  max-height: 200px;
}

.fade-height-enter-from,
.fade-height-leave-to {
  opacity: 0;
  transform: translateY(-10px);
  max-height: 0;
}
</style>
