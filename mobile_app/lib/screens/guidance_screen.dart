import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:erh_safety_app/providers/guidance_provider.dart';
import 'package:erh_safety_app/services/storage_service.dart';

class GuidanceScreen extends StatefulWidget {
  const GuidanceScreen({super.key});

  @override
  State<GuidanceScreen> createState() => _GuidanceScreenState();
}

class _GuidanceScreenState extends State<GuidanceScreen> {
  String? _currentZone;
  String? _targetZone;
  
  final List<String> _zones = ['Z1', 'Z2', 'Z3', 'Z4'];
  
  @override
  void initState() {
    super.initState();
    _loadCurrentZone();
  }
  
  void _loadCurrentZone() {
    final currentZone = StorageService.getCurrentZone();
    setState(() {
      _currentZone = currentZone ?? 'Z1';
      _targetZone = _currentZone;
    });
  }
  
  Future<void> _fetchGuidance() async {
    if (_currentZone == null || _targetZone == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select current and target zones')),
      );
      return;
    }
    
    final guidanceProvider = Provider.of<GuidanceProvider>(context, listen: false);
    await guidanceProvider.fetchGuidance(
      zoneId: _currentZone!,
      currentZone: _currentZone!,
      targetZone: _targetZone!,
    );
    
    if (!mounted) return;
    
    if (guidanceProvider.error != null) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(guidanceProvider.error!),
          backgroundColor: Colors.red,
        ),
      );
    }
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Guidance'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Get personalized guidance to your destination',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 24),
            DropdownButtonFormField<String>(
              value: _currentZone,
              decoration: const InputDecoration(
                labelText: 'Current Zone',
                border: OutlineInputBorder(),
              ),
              items: _zones.map((zone) {
                return DropdownMenuItem(
                  value: zone,
                  child: Text('Zone $zone'),
                );
              }).toList(),
              onChanged: (value) {
                setState(() {
                  _currentZone = value;
                });
                if (value != null) {
                  StorageService.saveCurrentZone(value);
                }
              },
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<String>(
              value: _targetZone,
              decoration: const InputDecoration(
                labelText: 'Target Zone',
                border: OutlineInputBorder(),
              ),
              items: _zones.map((zone) {
                return DropdownMenuItem(
                  value: zone,
                  child: Text('Zone $zone'),
                );
              }).toList(),
              onChanged: (value) {
                setState(() {
                  _targetZone = value;
                });
              },
            ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: _fetchGuidance,
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text('Get Guidance'),
            ),
            const SizedBox(height: 24),
            Consumer<GuidanceProvider>(
              builder: (context, guidanceProvider, child) {
                if (guidanceProvider.isLoading) {
                  return const Center(
                    child: Padding(
                      padding: EdgeInsets.all(32.0),
                      child: CircularProgressIndicator(),
                    ),
                  );
                }
                
                if (guidanceProvider.currentGuidance == null) {
                  return const Center(
                    child: Padding(
                      padding: EdgeInsets.all(32.0),
                      child: Text('Select zones and click "Get Guidance"'),
                    ),
                  );
                }
                
                return _buildGuidanceContent(guidanceProvider.currentGuidance!);
              },
            ),
          ],
        ),
      ),
    );
  }
  
  Widget _buildGuidanceContent(Map<String, dynamic> guidance) {
    final avoidZones = guidance['avoid_zones'] as List<dynamic>? ?? [];
    final recommendedPath = guidance['recommended_path'] as List<dynamic>? ?? [];
    final instructions = guidance['instructions'] as List<dynamic>? ?? [];
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (avoidZones.isNotEmpty) ...[
          Card(
            color: Colors.red.shade50,
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      const Icon(Icons.warning, color: Colors.red),
                      const SizedBox(width: 8),
                      Text(
                        'Avoid Zones',
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: Colors.red,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Wrap(
                    spacing: 8,
                    children: avoidZones.map((zone) {
                      return Chip(
                        label: Text('Zone $zone'),
                        backgroundColor: Colors.red.shade100,
                      );
                    }).toList(),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
        ],
        if (recommendedPath.isNotEmpty) ...[
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Recommended Path',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Wrap(
                    spacing: 8,
                    children: recommendedPath.asMap().entries.map((entry) {
                      final index = entry.key;
                      final zone = entry.value;
                      return Chip(
                        label: Text('${index + 1}. Zone $zone'),
                        avatar: index == 0
                            ? const Icon(Icons.location_on, size: 18)
                            : index == recommendedPath.length - 1
                                ? const Icon(Icons.flag, size: 18)
                                : const Icon(Icons.arrow_forward, size: 18),
                      );
                    }).toList(),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
        ],
        if (instructions.isNotEmpty) ...[
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Instructions',
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  ...instructions.map((instruction) {
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 8),
                      child: Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Icon(Icons.arrow_right, size: 20),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(instruction.toString()),
                          ),
                        ],
                      ),
                    );
                  }),
                ],
              ),
            ),
          ),
        ],
      ],
    );
  }
}

