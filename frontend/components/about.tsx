import { Button } from "@/components/ui/button"
import { CheckCircle } from "lucide-react"
import Link from "next/link"

const achievements = [
  "50+ successful projects delivered",
  "99.9% client satisfaction rate",
  "24/7 support and maintenance",
  "Award-winning design team",
]

export function About() {
  return (
    <section className="py-20">
      <div className="container mx-auto px-4">
        <div className="grid gap-12 lg:grid-cols-2 items-center">
          <div>
            <h2 className="text-3xl md:text-4xl font-bold mb-6">Why Choose Our Company?</h2>
            <p className="text-lg text-muted-foreground mb-8">
              With years of experience in the industry, we've helped hundreds of businesses transform their digital
              presence and achieve their goals. Our team of experts combines technical excellence with creative
              innovation to deliver solutions that make a real impact.
            </p>

            <div className="space-y-4 mb-8">
              {achievements.map((achievement, index) => (
                <div key={index} className="flex items-center space-x-3">
                  <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0" />
                  <span>{achievement}</span>
                </div>
              ))}
            </div>

            <Button size="lg" asChild>
              <Link href="/about">Learn More About Us</Link>
            </Button>
          </div>

          <div className="relative">
            <div className="aspect-square bg-gradient-to-br from-primary/20 to-purple-500/20 rounded-2xl p-8">
              <div className="w-full h-full bg-background rounded-xl shadow-2xl flex items-center justify-center">
                <div className="text-center">
                  <div className="text-4xl font-bold text-primary mb-2">5+</div>
                  <div className="text-muted-foreground">Years Experience</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
