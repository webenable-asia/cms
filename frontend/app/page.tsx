import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Navigation } from '@/components/navigation'
import { 
  Code2, 
  Smartphone, 
  Palette, 
  Zap, 
  Shield, 
  HeadphonesIcon,
  CheckCircle,
  Star,
  Award,
  Users,
  Clock,
  ArrowRight,
  Play,
  Mail,
  Phone,
  MapPin
} from 'lucide-react'

export default function Home() {
  const services = [
    {
      icon: Code2,
      title: "Web Development",
      description: "Custom web applications built with modern technologies like React, Next.js, and TypeScript."
    },
    {
      icon: Smartphone,
      title: "Mobile Apps",
      description: "Native and cross-platform mobile applications for iOS and Android devices."
    },
    {
      icon: Palette,
      title: "UI/UX Design",
      description: "Beautiful, intuitive designs that provide exceptional user experiences."
    },
    {
      icon: Zap,
      title: "Performance",
      description: "Lightning-fast applications optimized for speed and user engagement."
    },
    {
      icon: Shield,
      title: "Security",
      description: "Enterprise-grade security measures to protect your data and users."
    },
    {
      icon: HeadphonesIcon,
      title: "Support",
      description: "Ongoing support and maintenance to keep your applications running smoothly."
    }
  ]

  const stats = [
    {
      icon: CheckCircle,
      number: "50+",
      label: "successful projects delivered"
    },
    {
      icon: Star,
      number: "99.9%",
      label: "client satisfaction rate"
    },
    {
      icon: Clock,
      number: "24/7",
      label: "support and maintenance"
    },
    {
      icon: Award,
      number: "5+",
      label: "Years Experience"
    }
  ]

  return (
    <div className="min-h-screen bg-background">
      <Navigation />

      {/* Hero Section */}
      <section className="relative pt-16 pb-20 lg:pt-24 lg:pb-28">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-4xl mx-auto">
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold text-foreground mb-6 leading-tight">
              Build Amazing{" "}
              <span className="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                Digital Experiences
              </span>
            </h1>
            <p className="text-xl md:text-2xl text-muted-foreground mb-8 max-w-3xl mx-auto leading-relaxed">
              We help businesses create exceptional web applications, mobile apps, and digital 
              solutions that drive growth and engage customers.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <Link href="/contact">
                <Button size="lg" className="text-lg px-8 py-6 h-auto">
                  Get Started
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Button>
              </Link>
              <Button variant="outline" size="lg" className="text-lg px-8 py-6 h-auto">
                <Play className="mr-2 h-5 w-5" />
                Watch Demo
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Services Section */}
      <section className="py-20 bg-muted/50">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl lg:text-5xl font-bold text-foreground mb-4">
              Everything You Need to Succeed
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Our comprehensive suite of services helps you build, launch, and scale your 
              digital presence with confidence.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            {services.map((service, index) => {
              const Icon = service.icon
              return (
                <Card key={index} className="border-border hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
                  <CardHeader>
                    <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center mb-4">
                      <Icon className="h-6 w-6 text-primary" />
                    </div>
                    <CardTitle className="text-xl font-semibold text-card-foreground">
                      {service.title}
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-muted-foreground leading-relaxed">
                      {service.description}
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
            <h2 className="text-3xl md:text-4xl lg:text-5xl font-bold text-foreground mb-4">
              Why Choose Our Company?
            </h2>
            <p className="text-lg text-muted-foreground max-w-3xl mx-auto">
              With years of experience in the industry, we've helped hundreds of businesses 
              transform their digital presence and achieve their goals. Our team of experts 
              combines technical excellence with creative innovation to deliver solutions that 
              make a real impact.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, index) => {
              const Icon = stat.icon
              return (
                <div key={index} className="text-center">
                  <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
                    <Icon className="h-8 w-8 text-primary" />
                  </div>
                  <div className="text-3xl md:text-4xl font-bold text-foreground mb-2">
                    {stat.number}
                  </div>
                  <p className="text-muted-foreground">
                    {stat.label}
                  </p>
                </div>
              )
            })}
          </div>

          <div className="text-center mt-12">
            <Link href="/about">
              <Button variant="outline" size="lg" className="text-lg px-8 py-6 h-auto">
                Learn More About Us
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-primary text-primary-foreground">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-3xl md:text-4xl lg:text-5xl font-bold mb-4">
            Ready to Start Your Next Project?
          </h2>
          <p className="text-lg opacity-90 max-w-2xl mx-auto mb-8">
            Let's work together to bring your ideas to life and create something amazing. 
            Get in touch with our team today to discuss your project requirements.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link href="/contact">
              <Button variant="secondary" size="lg" className="text-lg px-8 py-6 h-auto">
                Start Your Project
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </Link>
            <Link href="/blog">
              <Button variant="outline" size="lg" className="text-lg px-8 py-6 h-auto">
                View Our Work
              </Button>
            </Link>
          </div>
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
