import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Users, Target, Award, Heart } from "lucide-react"

const values = [
  {
    icon: Users,
    title: "Team Collaboration",
    description: "We believe in the power of working together to achieve extraordinary results.",
  },
  {
    icon: Target,
    title: "Goal-Oriented",
    description: "Every project we undertake is driven by clear objectives and measurable outcomes.",
  },
  {
    icon: Award,
    title: "Excellence",
    description: "We strive for excellence in everything we do, from code quality to customer service.",
  },
  {
    icon: Heart,
    title: "Passion",
    description: "Our passion for technology and innovation drives us to create amazing solutions.",
  },
]

export default function AboutPage() {
  return (
    <div className="container mx-auto px-4 py-16">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-16">
          <h1 className="text-4xl font-bold mb-6">About Our Company</h1>
          <p className="text-xl text-muted-foreground leading-relaxed">
            We are a team of passionate developers and designers dedicated to creating exceptional digital experiences
            that help businesses grow and succeed in the modern world.
          </p>
        </div>

        <div className="grid gap-8 md:grid-cols-2 mb-16">
          <div>
            <h2 className="text-3xl font-bold mb-4">Our Story</h2>
            <p className="text-muted-foreground mb-4">
              Founded in 2020, our company started with a simple mission: to bridge the gap between innovative
              technology and practical business solutions. What began as a small team of developers has grown into a
              full-service digital agency.
            </p>
            <p className="text-muted-foreground">
              Today, we work with businesses of all sizes, from startups to enterprise companies, helping them leverage
              technology to achieve their goals and stay competitive in an ever-evolving digital landscape.
            </p>
          </div>
          <div>
            <h2 className="text-3xl font-bold mb-4">Our Mission</h2>
            <p className="text-muted-foreground mb-4">
              To empower businesses with cutting-edge technology solutions that drive growth, improve efficiency, and
              create meaningful connections with their customers.
            </p>
            <p className="text-muted-foreground">
              We believe that great technology should be accessible, reliable, and designed with the end user in mind.
              Every solution we create is built with these principles at its core.
            </p>
          </div>
        </div>

        <div className="mb-16">
          <h2 className="text-3xl font-bold text-center mb-8">Our Values</h2>
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
            {values.map((value, index) => (
              <Card key={index} className="text-center">
                <CardHeader>
                  <div className="mx-auto w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                    <value.icon className="w-6 h-6 text-primary" />
                  </div>
                  <CardTitle className="text-lg">{value.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <CardDescription>{value.description}</CardDescription>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        <div className="text-center bg-muted rounded-lg p-8">
          <h2 className="text-2xl font-bold mb-4">Ready to Work Together?</h2>
          <p className="text-muted-foreground mb-6">
            Let's discuss how we can help bring your ideas to life and achieve your business goals.
          </p>
          <a
            href="/contact"
            className="inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 bg-primary text-primary-foreground hover:bg-primary/90 h-10 px-4 py-2"
          >
            Get in Touch
          </a>
        </div>
      </div>
    </div>
  )
}
