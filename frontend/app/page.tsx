import { Hero } from "@/components/hero"
import { Features } from "@/components/features"
import { About } from "@/components/about"
import { CTA } from "@/components/cta"

export default function HomePage() {
  return (
    <main>
      <Hero />
      <Features />
      <About />
      <CTA />
    </main>
  )
}
