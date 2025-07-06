import { Button } from "@/components/ui/button"
import { ArrowRight, Play } from "lucide-react"
import Link from "next/link"

export function Hero() {
  return (
    <section className="relative py-20 lg:py-32 overflow-hidden">
      <div className="container mx-auto px-4">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-4xl md:text-6xl font-bold tracking-tight mb-6">
            Build Amazing Digital
            <span className="text-primary block">Experiences</span>
          </h1>
          <p className="text-xl text-muted-foreground mb-8 max-w-2xl mx-auto">
            We help businesses create exceptional web applications, mobile apps, and digital solutions that drive growth
            and engage customers.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" asChild>
              <Link href="/contact">
                Get Started
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button size="lg" variant="outline">
              <Play className="mr-2 h-4 w-4" />
              Watch Demo
            </Button>
          </div>
        </div>
      </div>

      {/* Background decoration */}
      <div className="absolute inset-0 -z-10 overflow-hidden">
        <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2">
          <div className="h-[600px] w-[600px] rounded-full bg-gradient-to-r from-primary/20 to-purple-500/20 blur-3xl" />
        </div>
      </div>
    </section>
  )
}
