import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Navigation } from '@/components/navigation'
import { 
  Users, 
  Target, 
  Award, 
  Heart,
  ArrowRight,
  Mail,
  Phone,
  MapPin,
  Lightbulb,
  Zap,
  Shield,
  Star
} from 'lucide-react'

export default function AboutPage() {
  const values = [
    {
      icon: Users,
      title: "Team Collaboration",
      description: "We believe in the power of working together to achieve extraordinary results."
    },
    {
      icon: Target,
      title: "Goal-Oriented",
      description: "Every project we undertake is driven by clear objectives and measurable outcomes."
    },
    {
      icon: Award,
      title: "Excellence",
      description: "We strive for excellence in everything we do, from code quality to customer service."
    },
    {
      icon: Heart,
      title: "Passion",
      description: "Our passion for technology and innovation drives us to create amazing solutions."
    }
  ]

  const stats = [
    { number: "50+", label: "Projects Completed" },
    { number: "25+", label: "Happy Clients" },
    { number: "5+", label: "Years Experience" },
    { number: "99%", label: "Client Satisfaction" }
  ]

  return (
    <div className="min-h-screen bg-background">
      <Navigation />

      {/* Hero Section */}
      <section className="pt-16 pb-20 lg:pt-24 lg:pb-28">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-4xl mx-auto">
            <h1 className="text-4xl md:text-6xl font-bold text-foreground mb-6">
              About Our Company
            </h1>
            <p className="text-xl md:text-2xl text-muted-foreground mb-8 max-w-3xl mx-auto leading-relaxed">
              We are a team of passionate developers and designers dedicated to creating 
              exceptional digital experiences that help businesses grow and succeed in the modern world.
            </p>
          </div>
        </div>
      </section>

      {/* Our Story Section */}
      <section className="py-20 bg-muted/50">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-4xl mx-auto">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-8 text-center">
              Our Story
            </h2>
            
            <div className="prose prose-lg mx-auto text-muted-foreground">
              <p className="text-lg leading-relaxed mb-6">
                Founded in 2020, our company started with a simple mission: to bridge the gap 
                between innovative technology and practical business solutions. What began as a 
                small team of developers has grown into a full-service digital agency.
              </p>
              
              <p className="text-lg leading-relaxed">
                Today, we work with businesses of all sizes, from startups to enterprise companies, 
                helping them leverage technology to achieve their goals and stay competitive in an 
                ever-evolving digital landscape.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Our Mission Section */}
      <section className="py-20">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-8">
              Our Mission
            </h2>
            
            <div className="bg-primary/10 rounded-2xl p-8 md:p-12">
              <p className="text-lg md:text-xl text-foreground leading-relaxed mb-6">
                To empower businesses with cutting-edge technology solutions that drive growth, 
                improve efficiency, and create meaningful connections with their customers.
              </p>
              
              <p className="text-lg text-muted-foreground">
                We believe that great technology should be accessible, reliable, and designed 
                with the end user in mind. Every solution we create is built with these principles 
                at its core.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Our Values Section */}
      <section className="py-20 bg-muted/50">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
              Our Values
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              These core values guide everything we do and shape how we work with our clients and each other.
            </p>
          </div>

          <div className="grid md:grid-cols-2 gap-8">
            {values.map((value, index) => {
              const Icon = value.icon
              return (
                <Card key={index} className="border-border hover:shadow-lg transition-all duration-300">
                  <CardHeader>
                    <div className="flex items-center mb-4">
                      <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mr-4">
                        <Icon className="h-6 w-6 text-primary" />
                      </div>
                      <CardTitle className="text-xl font-semibold text-card-foreground">
                        {value.title}
                      </CardTitle>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <p className="text-muted-foreground leading-relaxed">
                      {value.description}
                    </p>
                  </CardContent>
                </Card>
              )
            })}
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-20">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
              Our Impact
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Numbers that reflect our commitment to delivering exceptional results for our clients.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, index) => (
              <div key={index} className="text-center">
                <div className="text-4xl md:text-5xl font-bold text-primary mb-2">
                  {stat.number}
                </div>
                <p className="text-muted-foreground text-lg">
                  {stat.label}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-primary text-primary-foreground">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Ready to Work Together?
          </h2>
          <p className="text-lg opacity-90 max-w-2xl mx-auto mb-8">
            Let's discuss how we can help bring your ideas to life and achieve your business goals.
          </p>
          <Link href="/contact">
            <Button variant="secondary" size="lg" className="text-lg px-8 py-6 h-auto">
              Get in Touch
              <ArrowRight className="ml-2 h-5 w-5" />
            </Button>
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-secondary text-secondary-foreground py-16">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="lg:col-span-2">
              <h3 className="text-2xl font-bold mb-4">WebEnable</h3>
              <p className="text-muted-foreground mb-6 max-w-md">
                Building exceptional digital experiences that help businesses grow and succeed.
              </p>
              <div className="space-y-2">
                <div className="flex items-center text-muted-foreground">
                  <Mail className="h-4 w-4 mr-2" />
                  hello@webenable.asia
                </div>
                <div className="flex items-center text-muted-foreground">
                  <Phone className="h-4 w-4 mr-2" />
                  +66 (0) 123-456-789
                </div>
                <div className="flex items-center text-muted-foreground">
                  <MapPin className="h-4 w-4 mr-2" />
                  Bangkok, Thailand
                </div>
              </div>
            </div>
            
            <div>
              <h4 className="text-lg font-semibold mb-4">Company</h4>
              <ul className="space-y-2">
                <li><Link href="/about" className="text-muted-foreground hover:text-foreground transition-colors">About Us</Link></li>
                <li><Link href="/blog" className="text-muted-foreground hover:text-foreground transition-colors">Blog</Link></li>
                <li><Link href="/contact" className="text-muted-foreground hover:text-foreground transition-colors">Contact</Link></li>
              </ul>
            </div>
            
            <div>
              <h4 className="text-lg font-semibold mb-4">Services</h4>
              <ul className="space-y-2">
                <li className="text-muted-foreground">Web Development</li>
                <li className="text-muted-foreground">Mobile Apps</li>
                <li className="text-muted-foreground">UI/UX Design</li>
                <li className="text-muted-foreground">Consulting</li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-border mt-12 pt-8">
            <div className="flex flex-col md:flex-row justify-between items-center">
              <p className="text-muted-foreground mb-4 md:mb-0">
                © 2025 WebEnable. All rights reserved.
              </p>
              <div className="flex space-x-6">
                <Link href="#" className="text-muted-foreground hover:text-foreground transition-colors">Privacy Policy</Link>
                <Link href="#" className="text-muted-foreground hover:text-foreground transition-colors">Terms of Service</Link>
                <Link href="#" className="text-muted-foreground hover:text-foreground transition-colors">Cookie Policy</Link>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
