<template>
  <div class="flex items-center justify-center py-4 md:px-12 lg:px-24 font-sans whitespace-normal">
    <div class="w-full max-w-full bg-white overflow-hidden">
      <h1 class="text-3xl font-bold text-center text-gray-800 p-6 rounded-t-lg">
        Frequently Asked Questions
      </h1>
      <div class="p-6">
        <Accordion :items="faqItems" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Accordion, { AccordionItem } from "./Accordion.vue"

const faqItems: AccordionItem[] = [{
  title: "How does the AI specifically help with college research and shortlisting?",
  content: "Our AI platform utilizes advanced algorithms to analyze your academic transcripts, extracurricular activities, test scores, and personal preferences. It cross-references this data with an extensive database of universities, including their admission statistics, program offerings, faculty research, campus culture, and post-graduation outcomes. This allows us to generate a highly personalized shortlist of colleges that are not just a good fit academically, but also align with your social and career aspirations, saving you countless hours of manual research."
},
{
  title: "What kind of assistance does the AI provide for essays, personal statements, and supplements?",
  content: "The AI acts as your intelligent writing assistant. For essays, personal statements, and supplementary essays, it can help you brainstorm unique topics, outline compelling narratives, and refine your prose. It identifies areas for improvement in clarity, coherence, and impact, offering suggestions for stronger vocabulary, sentence structure, and thematic development. Crucially, it focuses on enhancing your authentic voice, ensuring the final submission truly reflects your individuality and experiences, rather than generating generic content."
},
{
  title: "Can the AI recommend specific courses or resources for standardized exams (e.g., SAT/ACT)?",
  content: "Absolutely. Based on your diagnostic test results and target scores, our AI platform identifies your strengths and weaknesses across different sections of standardized exams like the SAT or ACT. It then curates and recommends specific online courses, practice materials, and study strategies tailored to your learning style and areas needing improvement. This targeted approach helps optimize your study time and maximize your score potential."
},
{
  title: "How does the platform guide me on obtaining strong recommendation letters?",
  content: "Our AI provides strategic advice on selecting the most appropriate recommenders (teachers, counselors, mentors) who can best speak to your strengths and character. It offers guidance on how to approach them, what specific information or anecdotes to provide (e.g., a \"brag sheet\" template), and how to follow up professionally. The goal is to equip your recommenders with the insights they need to write personalized and impactful letters that highlight your unique contributions and potential."
},
{
  title: "What support does the AI offer for showcasing activities, honors, and applying for scholarships?",
  content: "The AI helps you craft compelling descriptions for your activities and honors, emphasizing their impact and your growth. It can suggest how to frame your experiences to align with the values of specific universities or scholarship programs. For scholarships, the platform identifies opportunities that match your academic profile, background, and interests from a vast database. It also provides tailored tips for essay prompts, application strategies, and helps you organize your materials to increase your chances of securing financial aid."
},
{
  title: "How does the AI ensure the privacy and security of my sensitive college application data?",
  content: "Data security and privacy are paramount. Our platform employs industry-standard encryption protocols (like SSL/TLS) for all data in transit and at rest. We utilize secure cloud infrastructure with robust access controls and regular security audits. Your personal information, academic records, and essay drafts are treated with the utmost confidentiality and are never shared with third parties without your explicit consent. We are committed to complying with relevant data protection regulations to safeguard your sensitive application details."
}]
</script>
