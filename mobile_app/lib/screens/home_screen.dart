import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:erh_safety_app/providers/cap_provider.dart';
import 'report_screen.dart';
import 'guidance_screen.dart';
import 'assistance_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;
  
  final List<Widget> _screens = [
    const HomeTab(),
    const GuidanceScreen(),
    const ReportScreen(),
    const AssistanceScreen(),
  ];
  
  @override
  void initState() {
    super.initState();
    // Fetch CAP messages for current zone
    WidgetsBinding.instance.addPostFrameCallback((_) {
      // TODO: Get current zone from location service
      Provider.of<CapProvider>(context, listen: false).fetchCapMessages('Z1');
    });
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('ERH Safety'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              // TODO: Navigate to settings
            },
          ),
        ],
      ),
      body: _screens[_selectedIndex],
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _selectedIndex,
        onTap: (index) {
          setState(() {
            _selectedIndex = index;
          });
        },
        type: BottomNavigationBarType.fixed,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(Icons.home),
            label: 'Home',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.directions),
            label: 'Guidance',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.report),
            label: 'Report',
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.help_outline),
            label: 'Assistance',
          ),
        ],
      ),
    );
  }
}

class HomeTab extends StatelessWidget {
  const HomeTab({super.key});
  
  @override
  Widget build(BuildContext context) {
    return Consumer<CapProvider>(
      builder: (context, capProvider, child) {
        if (capProvider.isLoading) {
          return const Center(child: CircularProgressIndicator());
        }
        
        if (capProvider.error != null) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  Icons.error_outline,
                  size: 64,
                  color: Colors.grey,
                ),
                const SizedBox(height: 16),
                Text(capProvider.error!),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () {
                    capProvider.fetchCapMessages('Z1'); // TODO: Get current zone
                  },
                  child: const Text('Retry'),
                ),
              ],
            ),
          );
        }
        
        if (capProvider.capMessages.isEmpty) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  Icons.check_circle_outline,
                  size: 64,
                  color: Colors.green,
                ),
                const SizedBox(height: 16),
                Text(
                  'No active alerts',
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const SizedBox(height: 8),
                Text(
                  'Your area is safe',
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: Colors.grey,
                  ),
                ),
              ],
            ),
          );
        }
        
        return ListView.builder(
          padding: const EdgeInsets.all(16),
          itemCount: capProvider.capMessages.length,
          itemBuilder: (context, index) {
            final message = capProvider.capMessages[index];
            return _buildCapMessageCard(context, message);
          },
        );
      },
    );
  }
  
  Widget _buildCapMessageCard(BuildContext context, Map<String, dynamic> message) {
    final info = message['info'] as List?;
    final headline = info?.isNotEmpty == true ? info![0]['headline'] : 'Alert';
    final description = info?.isNotEmpty == true ? info![0]['description'] : '';
    final severity = info?.isNotEmpty == true ? info![0]['severity'] : 'Unknown';
    
    Color severityColor;
    IconData severityIcon;
    switch (severity) {
      case 'Extreme':
        severityColor = Colors.red;
        severityIcon = Icons.warning;
        break;
      case 'Severe':
        severityColor = Colors.orange;
        severityIcon = Icons.warning_amber;
        break;
      default:
        severityColor = Colors.amber;
        severityIcon = Icons.info;
    }
    
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 4,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Icon(severityIcon, color: severityColor),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    headline,
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: severityColor,
                    ),
                  ),
                ),
              ],
            ),
            if (description.isNotEmpty) ...[
              const SizedBox(height: 8),
              Text(description),
            ],
          ],
        ),
      ),
    );
  }
}

