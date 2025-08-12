interface PricingCard {
  name: string
  priceMonthly: string
  priceYearly: string
  priceOneTime?: string
  description: string
  features: string[]
  highlight: boolean
}

export const pricing: PricingCard[] = [
  {
    name: "Student",
    priceMonthly: "$5",
    priceYearly: "$2.5",
    description: "Perfect for individual students",
    features: [
      "College search",
      "Application evaluation",
      "Personal roadmap builder"
    ],
    highlight: false
  },
  {
    name: "Team",
    priceMonthly: "$10",
    priceYearly: "$4.17",
    description: "For small groups & counselors",
    features: [
      "Everything in Starter",
      "Invite 10 students",
      "Exam preparation resources"
    ],
    highlight: true
  },
  {
    name: "Community",
    priceMonthly: "$30",
    priceYearly: "$15",
    priceOneTime: "$200",
    description: "Schools & large organizations",
    features: [
      "Everything in Pro",
      "Invite 40 students",
    ],
    highlight: false
  }
]