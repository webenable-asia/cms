import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Code, Smartphone, Palette, Zap, Shield, Users } from "lucide-react"

const features = [
  {
    icon: Code,
    title: "Web Development",
    description: "Custom web applications built with modern technologies like React, Next.js, and TypeScript.",
  },
  {
    icon: Smartphone,
    title: "Mobile Apps",
    description: "Native and cross-platform mobile applications for iOS and Android devices.",
  },
  {
    icon: Palette,
    title: "UI/UX Design",
    description: "Beautiful, intuitive designs that provide exceptional user experiences.",
  },
  {
    icon: Zap,
    title: "Performance",
    description: "Lightning-fast applications optimized for speed and user engagement.",
  },
  {
    icon: Shield,
    title: "Security",
    description: "Enterprise-grade security measures to protect your data and users.",
  },
  {
    icon: Users,
    title: "Support",
    description: "Ongoing support and maintenance to keep your applications running smoothly.",
  },
]

export function Features() {
  return (
    <section className="py-20 bg-muted/50">
      <div className="container mx-auto px-4">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">Everything You Need to Succeed</h2>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Our comprehensive suite of services helps you build, launch, and scale your digital presence with
            confidence.
          </p>
        </div>

        <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          {features.map((feature, index) => (
            <Card key={index} className="border-0 shadow-sm hover:shadow-md transition-shadow">
              <CardHeader>
                <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                  <feature.icon className="w-6 h-6 text-primary" />
                </div>
                <CardTitle className="text-xl">{feature.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">{feature.description}</CardDescription>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  )
}
